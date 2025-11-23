package servteams

import (
	"fmt"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/go-faster/errors"
)

type TeamsRepository interface {
	CreateTeamWithMembers(domain.TeamName, domain.Members) (domain.Team, error)
	GetMembersByTeamName(domain.TeamName) (domain.Members, error)
}

func (ts *TeamsService) NewTeam(team domain.Team) (domain.Team, error) {

	team, err := ts.repo.CreateTeamWithMembers(team.Name, team.Members)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicate) {
			return domain.Team{}, domain.ErrDuplicate
		}
		return domain.Team{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return team, nil
}

func (ts *TeamsService) TeamWithMembers(tName domain.TeamName) (domain.Team, error) {

	members, err := ts.repo.GetMembersByTeamName(tName)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Team{}, domain.ErrNotFound
		}
		return domain.Team{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	if members.Empty() {
		return domain.NewTeam(tName), domain.ErrNoContent
	}

	return domain.NewTeam(tName, members...), nil
}
