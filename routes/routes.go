package routes

import (
	"github.com/labstack/echo/v4"
	"po/internal/handlers/home"
)

func Register(e *echo.Echo) {
	e.GET("/", home.New().Home)
}
