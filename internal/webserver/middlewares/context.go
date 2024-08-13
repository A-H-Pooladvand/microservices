package middlewares

import (
	"github.com/labstack/echo/v4"
	"po/internal/app"
)

// Context wraps the echo.Context to the app.Context
func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return next(app.NewContext(ctx))
	}
}
