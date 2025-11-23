package restmembers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/members/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}
	return e
}

func TestRestMembers_UserSetIsActive(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		serviceSetup func(*mocks.MembersService, string)
		wantStatus   int
		wantErr      error
	}{
		{
			name: "successful update",
			requestBody: SetIsActiveRequest{
				UserID:   uuid.New().String(),
				IsActive: true,
			},
			serviceSetup: func(mockService *mocks.MembersService, userID string) {
				mockService.On(
					"SetMemberIsActive",
					mock.MatchedBy(func(m domain.Member) bool { return string(m.Id) == userID }),
				).Return(domain.Member{
					Id:     domain.MemberId(userID),
					Name:   "Test User",
					Status: domain.MemberStatusActive,
					Team:   domain.TeamName("backend"),
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "successful update to inactive",
			requestBody: SetIsActiveRequest{
				UserID:   uuid.New().String(),
				IsActive: false,
			},
			serviceSetup: func(mockService *mocks.MembersService, userID string) {
				mockService.On(
					"SetMemberIsActive",
					mock.MatchedBy(func(m domain.Member) bool { return string(m.Id) == userID }),
				).Return(domain.Member{
					Id:     domain.MemberId(userID),
					Name:   "Test User",
					Status: domain.MemberStatusInactive,
					Team:   domain.TeamName("backend"),
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "member not found",
			requestBody: SetIsActiveRequest{
				UserID:   uuid.New().String(),
				IsActive: true,
			},
			serviceSetup: func(mockService *mocks.MembersService, userID string) {
				mockService.On(
					"SetMemberIsActive",
					mock.MatchedBy(func(m domain.Member) bool {
						return string(m.Id) == userID
					}),
				).Return(domain.Member{}, domain.ErrNotFound)
			},
			wantErr: domain.HttpErrNotFound(),
		},
		{
			name:         "invalid request body",
			requestBody:  `{"UserID": "not-a-uuid", "IsActive": true}`,
			serviceSetup: func(mockService *mocks.MembersService, userID string) {},
			wantErr:      ErrBadReqBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewMembersService(t)

			var userID string
			if req, ok := tt.requestBody.(SetIsActiveRequest); ok {
				userID = req.UserID
			}
			tt.serviceSetup(mockService, userID)

			handler := New(mockService, zap.NewNop().Sugar())

			var bodyBytes []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = handler.UserSetIsActive(c)

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

			mockService.AssertExpectations(t)
		})
	}
}

func TestRestMembers_GetUserPeviewsById(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		serviceSetup func(*mocks.MembersService, string)
		wantStatus   int
		wantErr      error // может быть echo.HTTPError или domain.CustomHttpError
	}{
		{
			name:   "successful get",
			userID: uuid.New().String(),
			serviceSetup: func(mockService *mocks.MembersService, userID string) {
				mockService.On("MemberReviews", domain.MemberId(userID)).
					Return(domain.Member{Id: domain.MemberId(userID)}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "member not found",
			userID: uuid.New().String(),
			serviceSetup: func(mockService *mocks.MembersService, userID string) {
				mockService.On("MemberReviews", domain.MemberId(userID)).
					Return(domain.Member{}, domain.ErrNotFound)
			},
			wantErr: domain.HttpErrNotFound(),
		},
		{
			name:   "invalid UUID",
			userID: "not-a-uuid",
			serviceSetup: func(mockService *mocks.MembersService, userID string) {
			},
			wantErr: ErrBadReqParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewMembersService(t)
			tt.serviceSetup(mockService, tt.userID)

			handler := New(mockService, zap.NewNop().Sugar())

			req := httptest.NewRequest(http.MethodGet, "/users/getReview/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			err := handler.GetUserPeviewsById(c)

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
