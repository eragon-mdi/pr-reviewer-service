package srvrest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"github.com/labstack/echo/v4"
)

type RestSrv struct {
	*echo.Echo
	srv *http.Server
}

func New(cfg configs.Server) *RestSrv {
	e := echo.New()
	e.HideBanner = true

	return &RestSrv{
		Echo: e,
		srv: &http.Server{
			Addr:              fmt.Sprintf("%s:%s", cfg.AddressF, cfg.PortF),
			ReadTimeout:       cfg.ReadTimeoutF,
			WriteTimeout:      cfg.WriteTimeoutF,
			ReadHeaderTimeout: cfg.ReadHeaderTimeoutF,
			IdleTimeout:       cfg.IdleTimeoutF,
		},
	}
}

func (r *RestSrv) Serve() error {
	r.srv.Handler = r
	return r.srv.ListenAndServe()
}

func (r *RestSrv) Shutdown(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
