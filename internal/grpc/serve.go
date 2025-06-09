package grpc

import (
	"context"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/internal/handler"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/a-h-pooladvand/microservices/routes"
	"github.com/fatih/color"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func Invoke(lc fx.Lifecycle, config *config.Config, h *handler.Grpc) *grpc.Server {
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":"+config.GRPC.Addr)

			if err != nil {
				return err
			}

			// Register reflection service on gRPC server.
			reflection.Register(server)

			routes.RegisterGrpcRoutes(server, h)

			go func() {
				color.Green("gRPC server started on [::]:" + config.GRPC.Addr)

				if err = server.Serve(lis); err != nil {
					log.Error("GRPC serve error", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Stop()

			return nil
		},
	})

	return server
}
