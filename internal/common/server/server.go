package server

import (
	"context"
	"log"
	"time"

	srvrest "github.com/eragon-mdi/pr-reviewer-service/internal/common/api/rest"
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"golang.org/x/sync/errgroup"
)

const (
	ShutdownNoTimeout = -1
)

type Server interface {
	StartAll() error
	GracefulShutdown(timeoutSeconds int) error

	REST() *srvrest.RestSrv
}

type server struct {
	rest *srvrest.RestSrv
}

func New(cfg *configs.Servers) Server {
	return &server{
		rest: srvrest.New(cfg.REST),
	}
}

func (s *server) StartAll() error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("server.StartAll.recover: ", r)
		}
	}()

	eg := errgroup.Group{}

	eg.Go(func() error {
		return s.REST().Serve()
	})

	return eg.Wait()
}

func (s *server) REST() *srvrest.RestSrv {
	return s.rest
}

func (s *server) GracefulShutdown(timeoutSeconds int) error {
	ctx := context.Background()

	if timeoutSeconds >= 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		defer cancel()
	}

	if err := s.rest.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
