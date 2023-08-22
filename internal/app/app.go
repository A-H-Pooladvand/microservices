package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"po/pkg/log"
	"po/routes"
	"time"
)

func Serve(ctx context.Context) {
	e := echo.New()
	routes.Register(e)

	group, ctx := errgroup.WithContext(ctx)

	// Start the application
	group.Go(func() error {
		return e.Start(":8000")
	})

	// Graceful shutdown
	group.Go(func() error {
		// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		return e.Shutdown(ctx)
	})

	if err := group.Wait(); err != nil {
		log.Sugar.Error(err)
	}
}
