package webserver

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.elastic.co/apm/module/apmechov4/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"po/internal/webserver/middlewares"
)

func RegisterMiddlewares(e *echo.Echo, r *prometheus.Registry) {
	e.Use(middlewares.Context)
	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{Registerer: r}))
	e.Use(otelecho.Middleware("app"))
	e.Use(otelecho.Middleware("app"))
	e.Use(apmechov4.Middleware())
}
