package cmd

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"po/internal/app"
	"po/internal/grpc"
	"po/internal/providers"
	"po/internal/webserver"
	"po/pkg/zlog"
)

var cmd = &cobra.Command{
	Use:   "webserver",
	Short: "webserver",
	Long:  `Initializing...`,
}

func Execute() {
	// First it's mandatory to load our logger service
	// before any other thing
	zlog.Boot()

	// Todo:: Load generic environments from docker
	// Todo:: Load secrets from Vault service
	// Load .env file
	if err := godotenv.Load(); err != nil {
		zap.L().Panic("unable to load .env file", zap.Error(err))
	}

	ctx, cancel := app.WithCancel()
	defer cancel()

	defer panicRecover(cancel)

	// Boot third party services
	if err := providers.Boot(ctx); err != nil {
		zap.L().Panic("unable to boot service", zap.Error(err))
	}

	// Serve necessary protocols such as gRPC, HTTP etc...
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := serve(ctx); err != nil {
			zap.L().Panic("err to serve necessary protocols such as gRPC, HTTP etc", zap.Error(err))
		}
	}

	if err := cmd.Execute(); err != nil {
		zap.L().Panic("err", zap.Error(err))
	}
}

func serve(ctx app.Context) error {
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.Serve(ctx)
	})

	g.Go(func() error {
		return webserver.New().Serve(ctx)
	})

	return g.Wait()
}

func panicRecover(cancel context.CancelFunc) {
	if r := recover(); r != nil {
		cancel()
		zap.L().Error("panic recovery", zap.Any("panic message", r))
		app.Get().GracefulShutdown()
	}
}
