package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"po/internal/app"
	"po/internal/grpc"
	"po/internal/providers"
	"po/pkg/zlog"
)

var cmd = &cobra.Command{
	Use:   "app",
	Short: "app",
	Long:  `Initializing...`,
}

func Execute() {
	zlog.Boot()

	ctx, cancel := app.WithCancel()
	defer cancel()

	defer panicRecover(cancel)

	// Boot third party services
	if err := providers.Boot(ctx); err != nil {
		zlog.Panic(err)
	}

	// Serve necessary protocols such as gRPC, HTTP etc...
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := serve(ctx); err != nil {
			zlog.Panic(err)
		}
	}

	if err := cmd.Execute(); err != nil {
		zlog.Fatal(err)
	}
}

func serve(ctx app.Context) error {
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return grpc.Serve(ctx)
	})

	g.Go(func() error {
		return app.New().Serve(ctx)
	})

	return g.Wait()
}

func panicRecover(cancel context.CancelFunc) {
	if r := recover(); r != nil {
		fmt.Println("Hello World")
		cancel()
		zlog.Panic(r)
	}
}
