package main

import (
	"github.com/vysogota0399/gophermart/internal/api"
	"github.com/vysogota0399/gophermart/internal/config"
	"github.com/vysogota0399/gophermart/internal/grpc_server"
	"github.com/vysogota0399/gophermart/internal/grpc_server/handlers"
	"github.com/vysogota0399/gophermart/internal/grpc_server/repositories"
	"github.com/vysogota0399/gophermart/internal/logging"
	"github.com/vysogota0399/gophermart/internal/storage"
	services "github.com/vysogota0399/gophermart_protos/gen/services/denormalized_order"
	"go.uber.org/fx"
)

func main() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			logging.NewZapLogger,
			storage.LCNewStorage,
			grpc_server.NewServer,
			fx.Annotate(handlers.NewOrderDetailsHandler, fx.As(new(services.DenormalizedOrderServiceServer))),
			fx.Annotate(repositories.NewProductsRepository, fx.As(new(handlers.ProductsRepository))),

			api.NewHTTPServer,
			api.NewHandler,
		),
		fx.Supply(
			config.MustNewConfig(),
		),
		fx.Invoke(startHTTPServer, startGRPCServer, checkDBConnection),
	)
}

func startHTTPServer(*api.HTTPServer)     {}
func startGRPCServer(*grpc_server.Server) {}
func checkDBConnection(*storage.Storage)  {}
