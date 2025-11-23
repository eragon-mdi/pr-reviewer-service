package e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Health Check Tests
// ============================================================================

// TestHealthCheck_Success проверяет успешный ответ health check
func TestHealthCheck_Success(t *testing.T) {
	// Запрос: GET /health
	resp, err := HealthCheck()
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверка: статус 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверка: тело ответа содержит статус "ok"
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

// ============================================================================
// Teams Tests
// ============================================================================

// TestTeams_AddTeam_Success проверяет успешное создание команды
func TestTeams_AddTeam_Success(t *testing.T) {
	// Подготовка: создаем уникальное имя команды
	teamName := "e2e-team-success-" + uuid.New().String()[:8]
	userID := uuid.New().String()

	// Запрос: POST /teams/add
	req := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   userID,
				Username: "TestUser1",
				IsActive: true,
			},
		},
	}

	resp, err := AddTeam(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверка: статус 201 Created
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Проверка: команда создана с правильным именем
	var result AddTeamResponse
	err = ParseJSONResponse(resp, &result)
	require.NoError(t, err)
	assert.Equal(t, teamName, result.TeamName)
	assert.Len(t, result.Members, 1)
	assert.Equal(t, userID, result.Members[0].UserID)
}

// TestTeams_GetTeamByName_Success проверяет успешное получение команды по имени
func TestTeams_GetTeamByName_Success(t *testing.T) {
	// Подготовка: создаем команду
	teamName := "e2e-team-get-" + uuid.New().String()[:8]
	userID := uuid.New().String()

	createReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   userID,
				Username: "TestUser1",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание)
	resp1, err := AddTeam(createReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: GET /teams/get/:team_name
	resp2, err := GetTeamByName(teamName)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Проверка: статус 200 OK
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	// Проверка: команда найдена с правильным именем
	var result AddTeamResponse
	err = ParseJSONResponse(resp2, &result)
	require.NoError(t, err)
	assert.Equal(t, teamName, result.TeamName)
	assert.Len(t, result.Members, 1)
}

// TestTeams_GetTeamByName_NotFound проверяет ошибку при получении несуществующей команды
func TestTeams_GetTeamByName_NotFound(t *testing.T) {
	// Запрос: GET /teams/get/:team_name (несуществующая команда)
	nonExistentTeam := "e2e-team-not-found-" + uuid.New().String()[:8]
	resp, err := GetTeamByName(nonExistentTeam)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверка: статус 404 Not Found или 204 No Content (если команда не найдена)
	// В зависимости от реализации может быть 204 No Content для пустого результата
	assert.Contains(t, []int{http.StatusNotFound, http.StatusNoContent, http.StatusInternalServerError}, resp.StatusCode)

	// Если статус 404, проверяем код ошибки
	if resp.StatusCode == http.StatusNotFound {
		errResp, err := ParseErrorResponse(resp)
		require.NoError(t, err)
		assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
	}
}

// ============================================================================
// Users Tests
// ============================================================================

// TestUsers_SetIsActive_Success проверяет успешную установку активности пользователя
func TestUsers_SetIsActive_Success(t *testing.T) {
	// Подготовка: создаем команду с пользователем
	teamName := "e2e-team-setactive-" + uuid.New().String()[:8]
	userID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   userID,
				Username: "TestUser1",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: POST /users/setIsActive (установка is_active = false)
	req := SetIsActiveRequest{
		UserID:   userID,
		IsActive: false,
	}

	resp2, err := SetIsActive(req)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Проверка: статус 200 OK
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}

// TestUsers_SetIsActive_NotFound проверяет ошибку при установке активности несуществующего пользователя
func TestUsers_SetIsActive_NotFound(t *testing.T) {
	// Запрос: POST /users/setIsActive (несуществующий пользователь)
	req := SetIsActiveRequest{
		UserID:   uuid.New().String(),
		IsActive: false,
	}

	resp, err := SetIsActive(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Проверка: статус 404 Not Found
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Проверка: ошибка содержит код NOT_FOUND
	errResp, err := ParseErrorResponse(resp)
	require.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

// TestUsers_GetUserReviews_Success проверяет успешное получение ревью пользователя
func TestUsers_GetUserReviews_Success(t *testing.T) {
	// Подготовка: создаем команду с пользователями
	teamName := "e2e-team-reviews-" + uuid.New().String()[:8]
	authorID := uuid.New().String()
	reviewerID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   authorID,
				Username: "Author",
				IsActive: true,
			},
			{
				UserID:   reviewerID,
				Username: "Reviewer",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: POST /pullRequest/create (создание PR)
	prID := uuid.New().String()
	createPRReq := CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR",
		AuthorID:       authorID,
	}

	resp2, err := CreatePullRequest(createPRReq)
	require.NoError(t, err)
	resp2.Body.Close()
	require.Equal(t, http.StatusCreated, resp2.StatusCode)

	// Запрос 3: GET /users/getReview/:id (получение ревью)
	resp3, err := GetUserReviews(reviewerID)
	require.NoError(t, err)
	defer resp3.Body.Close()

	// Проверка: статус 200 OK
	assert.Equal(t, http.StatusOK, resp3.StatusCode)
}

// TestUsers_GetUserReviews_Empty проверяет получение пустого списка ревью
func TestUsers_GetUserReviews_Empty(t *testing.T) {
	// Подготовка: создаем команду с пользователем без PR
	teamName := "e2e-team-empty-reviews-" + uuid.New().String()[:8]
	userID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   userID,
				Username: "TestUser1",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: GET /users/getReview/:id (получение ревью для пользователя без PR)
	resp2, err := GetUserReviews(userID)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Проверка: статус 204 No Content или 200 OK с пустым списком
	assert.Contains(t, []int{http.StatusOK, http.StatusNoContent}, resp2.StatusCode)
}

// ============================================================================
// Pull Requests Tests
// ============================================================================

// TestPullRequests_Create_Success проверяет успешное создание PR
func TestPullRequests_Create_Success(t *testing.T) {
	// Подготовка: создаем команду с пользователями
	teamName := "e2e-team-pr-create-" + uuid.New().String()[:8]
	authorID := uuid.New().String()
	reviewerID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   authorID,
				Username: "Author",
				IsActive: true,
			},
			{
				UserID:   reviewerID,
				Username: "Reviewer",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: POST /pullRequest/create (создание PR)
	prID := uuid.New().String()
	req := CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR",
		AuthorID:       authorID,
	}

	resp2, err := CreatePullRequest(req)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Проверка: статус 201 Created
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)

	// Проверка: PR создан с правильными данными
	var result CreatePullRequestResponse
	err = ParseJSONResponse(resp2, &result)
	require.NoError(t, err)
	assert.Equal(t, prID, result.PR.PullRequestID)
	assert.Equal(t, "Test PR", result.PR.PullRequestName)
	assert.Equal(t, "OPEN", result.PR.Status)
}

// TestPullRequests_Create_Duplicate проверяет ошибку при создании дубликата PR
func TestPullRequests_Create_Duplicate(t *testing.T) {
	// Подготовка: создаем команду с пользователем
	teamName := "e2e-team-pr-dup-" + uuid.New().String()[:8]
	authorID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   authorID,
				Username: "Author",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: POST /pullRequest/create (создание PR)
	prID := uuid.New().String()
	req := CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR",
		AuthorID:       authorID,
	}

	resp2, err := CreatePullRequest(req)
	require.NoError(t, err)
	resp2.Body.Close()
	require.Equal(t, http.StatusCreated, resp2.StatusCode)

	// Запрос 3: POST /pullRequest/create (дубликат)
	resp3, err := CreatePullRequest(req)
	require.NoError(t, err)
	defer resp3.Body.Close()

	// Проверка: статус 409 Conflict
	assert.Equal(t, http.StatusConflict, resp3.StatusCode)

	// Проверка: ошибка содержит код PR_EXISTS или DUPLICATE
	errResp, err := ParseErrorResponse(resp3)
	require.NoError(t, err)
	assert.Contains(t, []string{"PR_EXISTS", "DUPLICATE"}, errResp.Error.Code)
}

// TestPullRequests_Merge_Success проверяет успешный мерж PR
func TestPullRequests_Merge_Success(t *testing.T) {
	// Подготовка: создаем команду с пользователями
	teamName := "e2e-team-pr-merge-" + uuid.New().String()[:8]
	authorID := uuid.New().String()
	reviewerID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   authorID,
				Username: "Author",
				IsActive: true,
			},
			{
				UserID:   reviewerID,
				Username: "Reviewer",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: POST /pullRequest/create (создание PR)
	prID := uuid.New().String()
	createPRReq := CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR",
		AuthorID:       authorID,
	}

	resp2, err := CreatePullRequest(createPRReq)
	require.NoError(t, err)
	resp2.Body.Close()
	require.Equal(t, http.StatusCreated, resp2.StatusCode)

	// Запрос 3: POST /pullRequest/merge (мерж PR)
	mergeReq := MergePullRequestRequest{
		PullRequestID: prID,
	}

	resp3, err := MergePullRequest(mergeReq)
	require.NoError(t, err)
	defer resp3.Body.Close()

	// Проверка: статус 200 OK
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

	// Проверка: PR имеет статус MERGED
	var result MergePullRequestResponse
	err = ParseJSONResponse(resp3, &result)
	require.NoError(t, err)
	assert.Equal(t, "MERGED", result.PR.Status)
}

// TestPullRequests_Merge_Idempotent проверяет идемпотентность мержа (повторный мерж)
func TestPullRequests_Merge_Idempotent(t *testing.T) {
	// Подготовка: создаем команду с пользователями
	teamName := "e2e-team-pr-merge-idem-" + uuid.New().String()[:8]
	authorID := uuid.New().String()
	reviewerID := uuid.New().String()

	createTeamReq := AddTeamRequest{
		TeamName: teamName,
		Members: []TeamMember{
			{
				UserID:   authorID,
				Username: "Author",
				IsActive: true,
			},
			{
				UserID:   reviewerID,
				Username: "Reviewer",
				IsActive: true,
			},
		},
	}

	// Запрос 1: POST /teams/add (создание команды)
	resp1, err := AddTeam(createTeamReq)
	require.NoError(t, err)
	resp1.Body.Close()
	require.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Запрос 2: POST /pullRequest/create (создание PR)
	prID := uuid.New().String()
	createPRReq := CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: "Test PR",
		AuthorID:       authorID,
	}

	resp2, err := CreatePullRequest(createPRReq)
	require.NoError(t, err)
	resp2.Body.Close()
	require.Equal(t, http.StatusCreated, resp2.StatusCode)

	// Запрос 3: POST /pullRequest/merge (первый мерж)
	mergeReq := MergePullRequestRequest{
		PullRequestID: prID,
	}

	resp3, err := MergePullRequest(mergeReq)
	require.NoError(t, err)
	resp3.Body.Close()
	require.Equal(t, http.StatusOK, resp3.StatusCode)

	// Запрос 4: POST /pullRequest/merge (повторный мерж - должен быть идемпотентным)
	resp4, err := MergePullRequest(mergeReq)
	require.NoError(t, err)
	defer resp4.Body.Close()

	// Проверка: статус 200 OK (не ошибка)
	assert.Equal(t, http.StatusOK, resp4.StatusCode)

	// Проверка: PR все еще имеет статус MERGED
	var result MergePullRequestResponse
	err = ParseJSONResponse(resp4, &result)
	require.NoError(t, err)
	assert.Equal(t, "MERGED", result.PR.Status)
}
