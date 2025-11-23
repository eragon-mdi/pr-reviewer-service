package resttransport

import (
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/api"
	restmembers "github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/members"
	restpullrequests "github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/pull-requests"
	restteams "github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest/teams"
	"go.uber.org/zap"
)

type RestTransport interface {
	api.TeamTransport
	api.UserTransport
	api.PullRequestTransport
}

type restTransport struct {
	*restteams.RestTeams
	*restmembers.RestMembers
	*restpullrequests.RestPullRequests
}

func New(s Service, l *zap.SugaredLogger) RestTransport {
	return &restTransport{
		RestTeams:        restteams.New(s, l),
		RestMembers:      restmembers.New(s, l),
		RestPullRequests: restpullrequests.New(s, l),
	}
}

type Service interface {
	restteams.TeamsService
	restmembers.MembersService
	restpullrequests.PullRequestService
}
