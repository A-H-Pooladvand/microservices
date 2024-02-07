package routes

import (
	"github.com/labstack/echo/v4"
	"po/internal/handlers"
)

func Register(e *echo.Echo) {
	e.GET("/", handlers.Home)
}
