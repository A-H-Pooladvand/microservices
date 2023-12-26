package app

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"po/cfg"
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

	if err := e.Start(":" + c.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("shutting down the server")
	}

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cc, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return e.Shutdown(cc)
}
