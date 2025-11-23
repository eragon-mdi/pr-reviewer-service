package service

import (
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	servmembers "github.com/eragon-mdi/pr-reviewer-service/internal/service/members"
	servpullrequests "github.com/eragon-mdi/pr-reviewer-service/internal/service/pull-requests"
	servteams "github.com/eragon-mdi/pr-reviewer-service/internal/service/teams"
	"github.com/eragon-mdi/pr-reviewer-service/internal/transport"
)

type service struct {
	*servteams.TeamsService
	*servmembers.MembersService
	*servpullrequests.PrService

	r   Repository
	cfg *configs.BussinesLogic
}

func New(r Repository, cfg *configs.BussinesLogic) transport.Service {
	ms := servmembers.NewMembersService(cfg, r)

	return &service{
		TeamsService:   servteams.NewTeamsService(r),
		MembersService: ms,
		PrService:      servpullrequests.NewPullRequestService(r, ms),

		r:   r,
		cfg: cfg,
	}
}

type Repository interface {
	servteams.Repository
	servmembers.Repository
	servpullrequests.Repository
}
