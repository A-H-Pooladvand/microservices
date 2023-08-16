package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Serve(ctx context.Context, port string) error {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return e.Start(":" + port)
}
