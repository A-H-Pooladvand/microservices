package metric

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"po/pkg/trace"
)

type RestHandler struct {
	service  *Service
	registry *prometheus.Registry
}

type RestHandlerParams struct {
	fx.In
	Service  *Service
	Tracer   trace.Tracer
	Registry *prometheus.Registry
}

func NewRestHandler(params RestHandlerParams) RestHandler {
	return RestHandler{
		service:  params.Service,
		registry: params.Registry,
	}
}

func (h RestHandler) Handle(c echo.Context) error {
	return echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: h.registry,
	})(c)
}
