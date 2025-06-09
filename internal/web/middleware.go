package web

import (
	"github.com/a-h-pooladvand/microservices/internal/web/middlewares"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmechov4/v2"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middlewares.Context)
	e.Use(echoprometheus.NewMiddleware("app"))
	e.Use(apmechov4.Middleware())
}
