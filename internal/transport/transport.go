package transport

import (
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/api"
	resttransport "github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest"
	"go.uber.org/zap"
)

type transport struct {
	resttransport.RestTransport
}

func New(s Service, l *zap.SugaredLogger) api.Transport {
	return &transport{
		RestTransport: resttransport.New(s, l),
	}
}

type Service interface {
	resttransport.Service
}
