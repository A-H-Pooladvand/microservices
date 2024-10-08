package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"po/internal/handlers"
)

// RegisterWebRoutes registers web routes.
func RegisterWebRoutes(e *echo.Echo, w *handlers.RestHandlers) {
	e.GET("/metrics", w.Metric.Handle)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("api/v1/")

	usersGroup := v1.Group("users")
	{
		usersGroup.GET("", w.User.Index)
		usersGroup.GET(":id", w.User.Show)
		usersGroup.POST("", w.User.Create)
	}
}
