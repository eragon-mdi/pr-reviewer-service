package main

import (
	"log"

	"github.com/eragon-mdi/pr-reviewer-service/internal/common/api"
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/logger"
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/server"
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/storage"
	"github.com/eragon-mdi/pr-reviewer-service/internal/repository"
	"github.com/eragon-mdi/pr-reviewer-service/internal/service"
	"github.com/eragon-mdi/pr-reviewer-service/internal/transport"
	resttransport "github.com/eragon-mdi/pr-reviewer-service/internal/transport/http/rest"

	rootctx "github.com/eragon-mdi/go-playground/server/root-ctx"
)

func main() {
	cfg := configs.MustLoad()

	l, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to set logger: %v", err)
	}

	rCtx, cancelAppCtx := rootctx.NotifyBackgroundCtxToShutdownSignal()
	defer cancelAppCtx()

	store, err := storage.Conn(rCtx, &cfg.Storages, storage.ConnTimeoutDefault)
	if err != nil {
		l.Error(err)
		return
	}

	r := repository.New(store)
	s := service.New(r, &cfg.BussinesLogic)
	t := transport.New(s, l)

	srv := server.New(&cfg.Servers)
	srv.REST().HTTPErrorHandler = resttransport.HTTPErrorHandler
	api.RegisterRoutes(srv, t, cfg.Servers.REST.HealthCheckRoute)
	go func() {
		if err := srv.StartAll(); err != nil {
			l.Errorf("failed to start servers: %v", err)
			cancelAppCtx()
		}
	}()

	<-rCtx.Done()

	if err := srv.GracefulShutdown(server.ShutdownNoTimeout); err != nil {
		l.Errorw("error during server shutdown", "cause", err)
	}
	if err := store.GracefulShutdown(); err != nil {
		l.Errorw("error disconnect store", "cause", err)
	}
}
