package webserver

import (
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmechov4/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"po/internal/webserver/middlewares"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middlewares.Context)
	e.Use(otelecho.Middleware("app"))
	e.Use(apmechov4.Middleware())
}
