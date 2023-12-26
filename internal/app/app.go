package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"po/cfg"
	"po/pkg/logger"
	"po/routes"
	"time"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Serve(ctx Context) error {
	c := cfg.NewApp()
	e := echo.New()
	e.HideBanner = true

	routes.Register(e)

	group, TODO := errgroup.WithContext(ctx)

	// Start the application
	group.Go(func() error {
		return e.Start(":" + c.Port)
	})

	// Graceful shutdown
	group.Go(func() error {
		// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
		// Use a buffered channel to avoid missing signals as recommended for signal.Notify
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		TODO, cancel := context.WithTimeout(TODO, 10*time.Second)
		defer cancel()

		return e.Shutdown(TODO)
	})

	return group.Wait()
}

// Recover Recovers from panic
func (a *App) Recover(cancel context.CancelFunc) {
	if r := recover(); r != nil {
		cancel()
		logger.Panic(r)
	}
}
