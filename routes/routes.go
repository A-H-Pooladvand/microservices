package routes

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func Register(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {

		zap.L().Info("Hello World")

		return c.JSON(http.StatusOK, "OK")
	})
}
