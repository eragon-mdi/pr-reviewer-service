package restpullrequests

import "go.uber.org/zap"

type RestPullRequests struct {
	s Service
	l *zap.SugaredLogger
}

func New(s Service, l *zap.SugaredLogger) *RestPullRequests {
	return &RestPullRequests{
		s: s,
		l: l,
	}
}

type Service interface {
	PullRequestService
}
