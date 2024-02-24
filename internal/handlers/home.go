package handlers

import (
	"github.com/labstack/echo/v4"
	"po/internal/app"
)

func Home(c echo.Context) error {
	ctx := c.(*app.Context)

	return ctx.R().Ok("Ok")
}
