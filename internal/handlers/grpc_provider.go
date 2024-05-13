package handlers

import (
	"go.uber.org/fx"
	"po/internal/handlers/user"
)

type GrpcHandlers struct {
	Ping Ping
	User user.GrpcHandler
}

type GrpcHandlerParams struct {
	fx.In
	//Ping Ping
	User user.GrpcHandler
}

func NewGrpcHandlers(params GrpcHandlerParams) *GrpcHandlers {
	return &GrpcHandlers{
		Ping: Ping{},
		User: params.User,
	}
}
