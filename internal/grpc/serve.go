package grpc

import (
	"context"
	"github.com/fatih/color"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"po/configs"
	"po/routes"
)

func Invoke(lc fx.Lifecycle, config *configs.App) *grpc.Server {
	server := grpc.NewServer()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":"+config.GrpcPort)

			if err != nil {
				return err
			}

			routes.RegisterGrpcRoutes(server)

			go func() {
				color.Green("gRPC server started on [::]:" + config.GrpcPort)

				if err = server.Serve(lis); err != nil {
					zap.L().Error("GRPC serve error", zap.Error(err))
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
