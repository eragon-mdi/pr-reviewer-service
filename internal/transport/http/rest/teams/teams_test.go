package restteams

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/teams/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupEcho() *echo.Echo {
	return echo.New()
}

func TestRestTeams_AddTeam(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		serviceSetup func(*mocks.TeamsService, TeamRequest)
		wantStatus   int
		wantErr      error
	}{
		{
			name: "successful create",
			requestBody: TeamRequest{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: uuid.New().String(), Username: "User1", IsActive: true},
				},
			},
			serviceSetup: func(mockService *mocks.TeamsService, req TeamRequest) {
				mockService.On(
					"NewTeam",
					mock.MatchedBy(func(team domain.Team) bool {
						return team.Name.String() == req.TeamName
					}),
				).Return(domain.NewTeam(
					domain.TeamName("backend"),
					domain.Member{
						Name:   "User1",
						Status: domain.MemberStatusActive,
					},
				), nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate team",
			requestBody: TeamRequest{
				TeamName: "backend",
				Members:  []TeamMember{},
			},
			serviceSetup: func(mockService *mocks.TeamsService, req TeamRequest) {
				mockService.On(
					"NewTeam",
					mock.MatchedBy(func(team domain.Team) bool {
						return team.Name.String() == req.TeamName
					}),
				).Return(domain.Team{}, domain.ErrDuplicate)
			},
			wantErr: domain.HttpErrTeamExists(),
		},
		{
			name:        "invalid request body",
			requestBody: "invalid json",
			serviceSetup: func(mockService *mocks.TeamsService, req TeamRequest) {
			},
			wantErr: ErrBadReqBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewTeamsService(t)

			var teamReq TeamRequest
			if req, ok := tt.requestBody.(TeamRequest); ok {
				teamReq = req
			}
			tt.serviceSetup(mockService, teamReq)

			handler := New(mockService, zap.NewNop().Sugar())

			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/teams/add", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.AddTeam(c)

			if tt.wantErr != nil {
				switch want := tt.wantErr.(type) {
				case *echo.HTTPError:
					var got *echo.HTTPError
					if errors.As(err, &got) {
						assert.Equal(t, want.Code, got.Code)
						assert.Equal(t, want.Message, got.Message)
					} else {
						t.Fatalf("expected echo.HTTPError, got %v", err)
					}
				case *domain.CustomHttpError:
					var got *domain.CustomHttpError
					if errors.As(err, &got) {
						assert.Equal(t, want.HttpCode, got.HttpCode)
						assert.Equal(t, want.Code, got.Code)
					} else {
						t.Fatalf("expected CustomHttpError, got %v", err)
					}
				default:
					t.Fatalf("unsupported wantErr type: %T", tt.wantErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
			}
		})
	}
}

func TestRestTeams_GetTeamByName(t *testing.T) {
	tests := []struct {
		name         string
		teamName     string
		serviceSetup func(*mocks.TeamsService, string)
		wantStatus   int
		wantErr      error
	}{
		{
			name:     "successful get",
			teamName: "backend",
			serviceSetup: func(mockService *mocks.TeamsService, teamName string) {
				mockService.On("TeamWithMembers", domain.TeamName(teamName)).
					Return(domain.NewTeam(
						domain.TeamName("backend"),
						domain.Member{
							Name:   "User1",
							Status: domain.MemberStatusActive,
						},
					), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "empty team name",
			teamName: "",
			serviceSetup: func(mockService *mocks.TeamsService, teamName string) {
			},
			wantErr: ErrBadReqParam,
		},
		{
			name:     "team not found",
			teamName: "nonexistent",
			serviceSetup: func(mockService *mocks.TeamsService, teamName string) {
				mockService.On("TeamWithMembers", domain.TeamName(teamName)).
					Return(domain.Team{}, domain.ErrNotFound)
			},
			wantErr: domain.HttpErrNotFound(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewTeamsService(t)
			tt.serviceSetup(mockService, tt.teamName)

			handler := New(mockService, zap.NewNop().Sugar())

			req := httptest.NewRequest(http.MethodGet, "/teams/get/:team_name", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("team_name")
			c.SetParamValues(tt.teamName)

			err := handler.GetTeamByName(c)

			if tt.wantErr != nil {
				switch want := tt.wantErr.(type) {
				case *echo.HTTPError:
					var got *echo.HTTPError
					if errors.As(err, &got) {
						assert.Equal(t, want.Code, got.Code)
						assert.Equal(t, want.Message, got.Message)
					} else {
						t.Fatalf("expected echo.HTTPError, got %v", err)
					}
				case *domain.CustomHttpError:
					var got *domain.CustomHttpError
					if errors.As(err, &got) {
						assert.Equal(t, want.HttpCode, got.HttpCode)
						assert.Equal(t, want.Code, got.Code)
					} else {
						t.Fatalf("expected CustomHttpError, got %v", err)
					}
				default:
					t.Fatalf("unsupported wantErr type: %T", tt.wantErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
			}
		})
	}
}
