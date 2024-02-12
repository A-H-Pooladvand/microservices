package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"po/internal/app"
)

func Home(c echo.Context) error {
	ctx := c.(*app.Context)

	zap.L().Error("Hello World")

	return ctx.R().Ok("Hello World")
}
