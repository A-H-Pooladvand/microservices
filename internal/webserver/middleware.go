package webserver

import (
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmechov4/v2"
	"po/internal/webserver/middlewares"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middlewares.Context)
	e.Use(apmechov4.Middleware())
}
