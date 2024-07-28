package handlers

import (
	"go.uber.org/fx"
	"po/internal/handlers/metric"
	"po/internal/handlers/user"
)

type RestHandlers struct {
	User   user.RestHandler
	Metric metric.RestHandler
}

type RestHandlerParams struct {
	fx.In
	User   user.RestHandler
	Metric metric.RestHandler
}

func NewRestHandlers(params RestHandlerParams) *RestHandlers {
	return &RestHandlers{
		User:   params.User,
		Metric: params.Metric,
	}
}
