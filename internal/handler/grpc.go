package handler

import (
	"github.com/a-h-pooladvand/microservices/internal/handler/grpc/user"
	"go.uber.org/fx"
)

type Grpc struct {
	Ping Ping
	User user.GrpcHandler
}

type GrpcParams struct {
	fx.In
	//Ping Ping
	User user.GrpcHandler
}

func NewGrpc(params GrpcParams) *Grpc {
	return &Grpc{
		Ping: Ping{},
		User: params.User,
	}
}
