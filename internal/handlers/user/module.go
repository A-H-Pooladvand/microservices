package user

import (
	"go.uber.org/fx"
)

// Module is the fx module for the user package
var Module = fx.Module(
	"user",
	fx.Provide(
		NewRestHandler,
		NewGrpcHandler,
		NewService,
		NewRepository,
	),
)
