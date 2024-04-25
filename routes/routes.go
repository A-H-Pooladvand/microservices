package routes

import (
	"github.com/labstack/echo/v4"
	"po/internal/handlers"
)

func RegisterWebRoutes(e *echo.Echo, w *handlers.WebHandlers) {
	usersGroup := e.Group("users")
	{
		usersGroup.GET("", w.User.Index)
	}
}
