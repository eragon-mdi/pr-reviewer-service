package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var baseURL string

// setBaseURL устанавливает baseURL для запросов
func setBaseURL(url string) {
	baseURL = url
}

// getBaseURL возвращает baseURL для запросов
func getBaseURL() string {
	return baseURL
}

// ============================================================================
// Health Check Requests
// ============================================================================

// HealthCheck выполняет GET запрос к /health
func HealthCheck() (*http.Response, error) {
	return http.Get(baseURL + "/health")
}

// ============================================================================
// Teams Requests
// ============================================================================

// AddTeamRequest представляет запрос на создание команды
type AddTeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// TeamMember представляет участника команды
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// AddTeamResponse представляет ответ на создание команды
type AddTeamResponse struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// AddTeam выполняет POST запрос к /teams/add
func AddTeam(req AddTeamRequest) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/teams/add", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(httpReq)
}

// GetTeamByName выполняет GET запрос к /teams/get/:team_name
func GetTeamByName(teamName string) (*http.Response, error) {
	return http.Get(baseURL + "/teams/get/" + teamName)
}

// ============================================================================
// Users Requests
// ============================================================================

// SetIsActiveRequest представляет запрос на установку активности пользователя
type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// SetIsActive выполняет POST запрос к /users/setIsActive
func SetIsActive(req SetIsActiveRequest) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/users/setIsActive", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(httpReq)
}

// GetUserReviews выполняет GET запрос к /users/getReview/:id
func GetUserReviews(userID string) (*http.Response, error) {
	return http.Get(baseURL + "/users/getReview/" + userID)
}

// ============================================================================
// Pull Requests Requests
// ============================================================================

// CreatePullRequestRequest представляет запрос на создание PR
type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

// CreatePullRequestResponse представляет ответ на создание PR
type CreatePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

// PullRequest представляет PR
type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

// CreatePullRequest выполняет POST запрос к /pullRequest/create
func CreatePullRequest(req CreatePullRequestRequest) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/pullRequest/create", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(httpReq)
}

// MergePullRequestRequest представляет запрос на мерж PR
type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

// MergePullRequestResponse представляет ответ на мерж PR
type MergePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

// MergePullRequest выполняет POST запрос к /pullRequest/merge
func MergePullRequest(req MergePullRequestRequest) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/pullRequest/merge", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(httpReq)
}

// ReassignUserForPullRequestRequest представляет запрос на переназначение ревьювера
type ReassignUserForPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_reviewer_id"`
}

// ReassignUserForPullRequestResponse представляет ответ на переназначение ревьювера
type ReassignUserForPullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

// ReassignUserForPullRequest выполняет POST запрос к /pullRequest/reassign
func ReassignUserForPullRequest(req ReassignUserForPullRequestRequest) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/pullRequest/reassign", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(httpReq)
}

// ============================================================================
// Helper Functions
// ============================================================================

// ParseJSONResponse парсит JSON ответ в указанную структуру
// ВАЖНО: не закрывает resp.Body, вызывающий код должен закрыть его сам
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	return json.Unmarshal(body, v)
}

// ParseErrorResponse парсит ошибку из ответа
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ParseErrorResponse парсит ошибку из ответа
// ВАЖНО: не закрывает resp.Body, вызывающий код должен закрыть его сам
func ParseErrorResponse(resp *http.Response) (*ErrorResponse, error) {
	var errResp ErrorResponse
	if err := ParseJSONResponse(resp, &errResp); err != nil {
		return nil, err
	}
	return &errResp, nil
}
