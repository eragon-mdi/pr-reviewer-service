package restteams

import (
	"errors"
	"net/http"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/eragon-mdi/pr-reviewer-service/pkg/validator"
	"github.com/labstack/echo/v4"
)

var (
	ErrBadReqParam = echo.NewHTTPError(http.StatusBadRequest, "bad req param")
	ErrBadReqBody  = echo.NewHTTPError(http.StatusBadRequest, "bad req body")
)

type TeamsService interface {
	NewTeam(team domain.Team) (domain.Team, error)
	TeamWithMembers(tName domain.TeamName) (domain.Team, error)
}

func (ts *RestTeams) AddTeam(c echo.Context) error {
	var TeamReq = &TeamRequest{}

	l := ts.l.With("req", TeamReq)
	l.Infof("AddTeam called")

	if err := c.Bind(TeamReq); err != nil {
		l.Errorf("failed to bind request: %v", err)
		return ErrBadReqBody
	}

	if err := validate(c, TeamReq); err != nil {
		l.Errorf("failed validate: %v", err)
		return ErrBadReqBody
	}

	newTeam, err := ts.s.NewTeam(TeamReq.domain())
	if err != nil {
		l.Errorf("failed create team: %v", err)

		if errors.Is(err, domain.ErrDuplicate) {
			return domain.HttpErrTeamExists()
		}
		return domain.ErrInternal
	}

	l = l.With("id", newTeam)
	l.Infof("team created successfully")
	return c.JSON(http.StatusCreated, teamResponse(newTeam))
}

func (ts *RestTeams) GetTeamByName(c echo.Context) error {
	tName := c.Param("team_name")

	l := ts.l.With("team_name", tName)
	l.Infof("GetTeamByName called")

	if tName == "" {
		l.Errorf("team_name is empty")
		return ErrBadReqParam
	}

	team, err := ts.s.TeamWithMembers(domain.TeamName(tName))
	if err != nil {
		l.Errorf("failed get team: %v", err)

		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrNoContent) {
			return domain.HttpErrNotFound()
		}
		return domain.ErrInternal
	}

	l = l.With("team", team.Name.String())
	l.Infof("team fetched successfully")

	return c.JSON(http.StatusOK, teamResponse(team))
}

func validate(c echo.Context, structure any) error {
	return validator.Validate(c.Request().Context(), structure)
}
