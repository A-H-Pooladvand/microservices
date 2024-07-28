package metric

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewRestHandler,
		NewService,
		NewRepository,
	),
)
