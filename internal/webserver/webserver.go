package webserver

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"net/http"
	"os"
	"os/signal"
	"po/configs"
	"po/internal/handlers"
	"po/routes"
)

func Invoke(lc fx.Lifecycle, c *configs.App, w *handlers.RestHandlers) *echo.Echo {
	e := echo.New()

	RegisterMiddlewares(e)
	e.HideBanner = true
	e.HidePort = true
	routes.RegisterWebRoutes(e, w)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			color.Green("⇨ http server started on http://127.0.0.1:%v\n", c.AppPort)

			go func() {
				if err := e.Start(":" + c.AppPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
					e.Logger.Fatal("shutting down the server")
				}

				// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
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
