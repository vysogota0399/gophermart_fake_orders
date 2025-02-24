package api

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/vysogota0399/gophermart/internal/config"
	"go.uber.org/fx"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, h *Handler) *HTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/api/order", h)

	s := &HTTPServer{
		srv: &http.Server{Addr: cfg.HTTPAddress, Handler: mux, ReadHeaderTimeout: time.Minute},
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		},
	)

	return s
}

func (s *HTTPServer) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}

	go s.srv.Serve(ln)

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
