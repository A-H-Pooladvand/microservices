package middlewares

import (
	"github.com/labstack/echo/v4"
	"po/internal/app"
)

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cc := &app.Context{
			Context: ctx,
		}

		return next(cc)
	}
}
