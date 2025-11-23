package sqlrepo

import (
	sqlstore "github.com/eragon-mdi/go-playground/storage/sql"
	servmembers "github.com/eragon-mdi/pr-reviewer-service/internal/service/members"
	servpullrequests "github.com/eragon-mdi/pr-reviewer-service/internal/service/pull-requests"
	servteams "github.com/eragon-mdi/pr-reviewer-service/internal/service/teams"
)

type SqlRepo interface {
	servteams.Repository
	servmembers.Repository
	servpullrequests.Repository
}

type sqlRepo struct {
	*teamsRepo
	*membersRepo
	*pullRequestsRepo
}

func New(s sqlstore.Storage) SqlRepo {
	return &sqlRepo{
		teamsRepo:        NewTeamsRepo(s),
		membersRepo:      NewMembersRepo(s),
		pullRequestsRepo: NewPullRequestsRepo(s),
	}
}
