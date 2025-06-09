package web

import (
	"context"
	"errors"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/internal/handler"
	"github.com/a-h-pooladvand/microservices/pkg/validator"
	"github.com/a-h-pooladvand/microservices/routes"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"net/http"
	"os"
	"os/signal"
)

func Invoke(
	lc fx.Lifecycle,
	c *config.Config,
	w *handler.Http,
) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Validator = validator.New()

	RegisterMiddlewares(e)

	routes.RegisterWebRoutes(e, w)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			color.Green("⇨ http server started on http://127.0.0.1:%v\n", c.App.Port)

			go func() {
				if err := e.Start(":" + c.App.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
					e.Logger.Fatal("shutting down the server")
				}

				// WaitDuration for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
				// Use a buffered channel to avoid missing signals as recommended for signal.Notify
				quit := make(chan os.Signal, 1)
				signal.Notify(quit, os.Interrupt)
				<-quit
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
