package handlers

import (
	"go.uber.org/fx"
	"po/internal/handlers/user"
)

type RestHandlers struct {
	User user.RestHandler
}

type RestHandlerParams struct {
	fx.In
	User user.RestHandler
}

func NewRestHandlers(params RestHandlerParams) *RestHandlers {
	return &RestHandlers{
		User: params.User,
	}
}
