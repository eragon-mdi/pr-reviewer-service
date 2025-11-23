package restmembers

import "go.uber.org/zap"

type RestMembers struct {
	s Service
	l *zap.SugaredLogger
}

func New(s Service, l *zap.SugaredLogger) *RestMembers {
	return &RestMembers{
		s: s,
		l: l,
	}
}

type Service interface {
	MembersService
}
