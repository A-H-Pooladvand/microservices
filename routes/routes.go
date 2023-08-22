package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Register(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})
}
