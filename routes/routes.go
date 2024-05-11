package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"po/internal/handlers"
)

func RegisterWebRoutes(e *echo.Echo, w *handlers.WebHandlers) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	usersGroup := e.Group("users")
	{
		usersGroup.GET("", w.User.Index)
	}
}
