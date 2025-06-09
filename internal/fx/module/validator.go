package module

import (
	"github.com/a-h-pooladvand/microservices/pkg/validator"
	"go.uber.org/fx"
)

var Validator = fx.Module(
	"validator",
	fx.Provide(
		validator.New,
	),
)
