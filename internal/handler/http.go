package handler

import (
	"github.com/a-h-pooladvand/microservices/internal/handler/http/user"
	"go.uber.org/fx"
)

type Http struct {
	User user.HttpHandler
}

type HttpParams struct {
	fx.In
	User user.HttpHandler
}

func NewHttp(params HttpParams) *Http {
	return &Http{
		User: params.User,
	}
}
