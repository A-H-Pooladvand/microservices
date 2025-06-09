package fx

import (
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/internal/fx/invoke"
	"github.com/a-h-pooladvand/microservices/internal/fx/module"
	"github.com/a-h-pooladvand/microservices/internal/handler"
	"github.com/a-h-pooladvand/microservices/internal/web"
	"go.uber.org/fx"
	"time"
)

func App() *fx.App {
	return fx.New(
		Options()...,
	)
}

func Options() []fx.Option {
	return []fx.Option{
		module.Redis,
		module.Jaeger,
		module.Vault,
		module.Validator,
		module.Prometheus,
		module.User,
		module.Postgres,

		// Invoke executes the given functions sequentially (in order).
		fx.Invoke(
			invoke.Log,
			invoke.Timezone,
			invoke.Faker,
			web.Invoke,
		),

		fx.Provide(
			config.New,
			handler.NewHttp,
		),
		fx.StopTimeout(1 * time.Minute),
	}
}
