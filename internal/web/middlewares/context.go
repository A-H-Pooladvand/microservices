package middlewares

import (
	"github.com/a-h-pooladvand/microservices/internal/app"
	"github.com/labstack/echo/v4"
)

// Context wraps the echo.Context to the app.Context
func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return next(app.NewContext(ctx))
	}
}
