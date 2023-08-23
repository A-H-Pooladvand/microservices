package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Register(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
