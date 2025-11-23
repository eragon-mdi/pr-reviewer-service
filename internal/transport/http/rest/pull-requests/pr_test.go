package restpullrequests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	restpullrequests "github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/pull-requests"
	"github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/pull-requests/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupEcho() *echo.Echo {
	return echo.New()
}

func TestRestPullRequests_CreatePullRequest(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		serviceSetup func(*mocks.PullRequestService, restpullrequests.CreatePRRequest)
		wantStatus   int
		wantErr      error
	}{
		{
			name: "successful create",
			requestBody: restpullrequests.CreatePRRequest{
				PullRequestID:   uuid.New().String(),
				PullRequestName: "Test PR",
				AuthorID:        uuid.New().String(),
			},
			serviceSetup: func(mockService *mocks.PullRequestService, req restpullrequests.CreatePRRequest) {
				mockService.On(
					"NewPullRequest",
					mock.MatchedBy(func(pr domain.PullRequestShort) bool { return true }),
				).Return(domain.PullRequest{
					Id:        domain.PrId(req.PullRequestID),
					Name:      domain.PrName(req.PullRequestName),
					AuthorId:  domain.MemberId(req.AuthorID),
					Status:    domain.PrStatusOpen,
					CreatedAt: time.Now(),
				}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate pr",
			requestBody: restpullrequests.CreatePRRequest{
				PullRequestID:   uuid.New().String(),
				PullRequestName: "Test PR",
				AuthorID:        uuid.New().String(),
			},
			serviceSetup: func(mockService *mocks.PullRequestService, req restpullrequests.CreatePRRequest) {
				mockService.On(
					"NewPullRequest",
					mock.MatchedBy(func(pr domain.PullRequestShort) bool { return true }),
				).Return(domain.PullRequest{}, domain.ErrDuplicate)
			},
			wantErr: domain.HttpErrPRExists(),
		},
		{
			name:        "invalid request body",
			requestBody: "invalid json",
			serviceSetup: func(mockService *mocks.PullRequestService, req restpullrequests.CreatePRRequest) {
			},
			wantErr: restpullrequests.ErrBadReqBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewPullRequestService(t)

			var createReq restpullrequests.CreatePRRequest
			if req, ok := tt.requestBody.(restpullrequests.CreatePRRequest); ok {
				createReq = req
			}

			tt.serviceSetup(mockService, createReq)

			handler := restpullrequests.New(mockService, zap.NewNop().Sugar())

			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreatePullRequest(c)

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

func TestRestPullRequests_MergePullRequest(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		serviceSetup func(*mocks.PullRequestService, string)
		wantStatus   int
		wantErr      error
	}{
		{
			name: "successful merge",
			requestBody: restpullrequests.MergePRRequest{
				PullRequestID: uuid.New().String(),
			},
			serviceSetup: func(mockService *mocks.PullRequestService, prID string) {
				mockService.On("Merge", domain.PrId(prID)).
					Return(domain.PullRequest{
						Id:       domain.PrId(prID),
						Status:   domain.PrStatusMerged,
						MergedAt: time.Now(),
					}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "pr not found",
			requestBody: restpullrequests.MergePRRequest{
				PullRequestID: uuid.New().String(),
			},
			serviceSetup: func(mockService *mocks.PullRequestService, prID string) {
				mockService.On("Merge", domain.PrId(prID)).
					Return(domain.PullRequest{}, domain.ErrNotFound)
			},
			wantErr: domain.HttpErrNotFound(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewPullRequestService(t)

			var prID string
			if req, ok := tt.requestBody.(restpullrequests.MergePRRequest); ok {
				prID = req.PullRequestID
			}
			tt.serviceSetup(mockService, prID)

			handler := restpullrequests.New(mockService, zap.NewNop().Sugar())

			bodyBytes, _ := json.Marshal(tt.requestBody)

			req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.MergePullRequest(c)

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

func TestRestPullRequests_ReassignUserForPullRequest(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		serviceSetup func(*mocks.PullRequestService, restpullrequests.ReassignPRRequest)
		wantStatus   int
		wantErr      error
	}{
		{
			name: "successful reassign",
			requestBody: restpullrequests.ReassignPRRequest{
				PullRequestID: uuid.New().String(),
				OldUserID:     uuid.New().String(),
			},
			serviceSetup: func(mockService *mocks.PullRequestService, req restpullrequests.ReassignPRRequest) {
				mockService.On(
					"Reasign",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					mock.MatchedBy(func(prReasMem domain.PrReasignMember) bool {
						return string(prReasMem.PrId) == req.PullRequestID &&
							string(prReasMem.MemberId) == req.OldUserID
					}),
				).Return(domain.PrWithReasignMember{
					PullRequest: domain.PullRequest{
						Id:     domain.PrId(req.PullRequestID),
						Status: domain.PrStatusOpen,
					},
					MemberId: domain.MemberId(uuid.New().String()),
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "pr not found",
			requestBody: restpullrequests.ReassignPRRequest{
				PullRequestID: uuid.New().String(),
				OldUserID:     uuid.New().String(),
			},
			serviceSetup: func(mockService *mocks.PullRequestService, req restpullrequests.ReassignPRRequest) {
				mockService.On(
					"Reasign",
					mock.MatchedBy(func(ctx context.Context) bool { return true }),
					mock.MatchedBy(func(prReasMem domain.PrReasignMember) bool {
						return string(prReasMem.PrId) == req.PullRequestID &&
							string(prReasMem.MemberId) == req.OldUserID
					}),
				).Return(domain.PrWithReasignMember{}, domain.ErrNotFound)
			},
			wantErr: domain.HttpErrNotFound(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			mockService := mocks.NewPullRequestService(t)

			var reassignReq restpullrequests.ReassignPRRequest
			if req, ok := tt.requestBody.(restpullrequests.ReassignPRRequest); ok {
				reassignReq = req
			}
			tt.serviceSetup(mockService, reassignReq)

			handler := restpullrequests.New(mockService, zap.NewNop().Sugar())

			bodyBytes, _ := json.Marshal(tt.requestBody)

			req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.ReassignUserForPullRequest(c)

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
