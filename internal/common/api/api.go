package api

import (
	"net/http"

	"github.com/eragon-mdi/pr-reviewer-service/internal/common/server"
	"github.com/labstack/echo/v4"
)

type Transport interface {
	TeamTransport
	UserTransport
	PullRequestTransport
}

type TeamTransport interface {
	AddTeam(echo.Context) error
	GetTeamByName(echo.Context) error
}

type UserTransport interface {
	UserSetIsActive(echo.Context) error
	GetUserPeviewsById(echo.Context) error
}

type PullRequestTransport interface {
	CreatePullRequest(echo.Context) error
	MergePullRequest(echo.Context) error
	ReassignUserForPullRequest(echo.Context) error
}

func RegisterRoutes(s server.Server, t Transport, healthCheckRoute string) {
	s.REST().GET(healthCheckRoute, healthCheck)

	teams := s.REST().Group("/teams")
	teams.POST("/add", t.AddTeam)
	teams.GET("/get/:team_name", t.GetTeamByName)

	users := s.REST().Group("/users")
	users.POST("/setIsActive", t.UserSetIsActive)
	users.GET("/getReview/:id", t.GetUserPeviewsById)

	pullRequest := s.REST().Group("/pullRequest")
	pullRequest.POST("/create", t.CreatePullRequest)
	pullRequest.POST("/merge", t.MergePullRequest)
	pullRequest.POST("/reassign", t.ReassignUserForPullRequest)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
