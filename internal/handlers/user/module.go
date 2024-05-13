package user

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewRestHandler,
		NewGrpcHandler,
		NewService,
		NewRepository,
	),
)
