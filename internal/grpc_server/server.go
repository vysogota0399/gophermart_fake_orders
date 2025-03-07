package grpc_server

import (
	"context"
	"net"

	"github.com/vysogota0399/gophermart/internal/config"
	"github.com/vysogota0399/gophermart_protos/gen/queries/order_details"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	handler order_details.QueryOrderDetailsServer
	cfg     *config.Config
	srv     *grpc.Server
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cfg.GRPCAddress)
	if err != nil {
		return err
	}

	order_details.RegisterQueryOrderDetailsServer(s.srv, s.handler)
	go s.srv.Serve(lis)

	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

func NewServer(h order_details.QueryOrderDetailsServer, lc fx.Lifecycle, cfg *config.Config) *Server {
	srv := &Server{cfg: cfg, srv: grpc.NewServer(), handler: h}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return srv.Start()
			},
			OnStop: func(ctx context.Context) error {
				srv.Stop()
				return nil
			},
		},
	)

	return srv
}
