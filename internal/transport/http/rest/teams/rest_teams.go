package restteams

import "go.uber.org/zap"

type RestTeams struct {
	s Service
	l *zap.SugaredLogger
}

func New(s Service, l *zap.SugaredLogger) *RestTeams {
	return &RestTeams{
		s: s,
		l: l,
	}
}

type Service interface {
	TeamsService
}
