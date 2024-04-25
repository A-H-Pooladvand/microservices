package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"po/configs"
	"po/internal/apm"
	"po/internal/app"
	"po/internal/db"
	"po/internal/grpc"
	"po/internal/handlers"
	"po/internal/handlers/user"
	"po/internal/vault"
	"po/internal/webserver"
	"po/pkg/cache"
	"po/pkg/log"
	"po/pkg/logstash"
	"po/pkg/rabbitmq"
	"po/pkg/redis"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the application",
	// Invoke the necessary protocols such as gRPC, HTTP, etc...
	Run: runApplication,
}

var DBModule = fx.Module(
	"db",
	fx.Provide(
		configs.NewPostgres,
		db.New,
	),
)

func runApplication(cmd *cobra.Command, args []string) {
	app.LoadEnvironmentVariablesInLocalEnv()

	application := fx.New(
		DBModule,
		user.Module,
		fx.Provide(
			// Loading configs
			configs.NewApp,
			configs.NewLogstash,
			configs.NewRabbitMQ,
			configs.NewRedis,
			vault.NewConfig,
			// Loading services
			logstash.New,
			vault.Provide,
			rabbitmq.Provide,
			handlers.NewWebHandlers,
			fx.Annotate(
				redis.Provide,
				fx.As(new(cache.Cache)),
			),
		),

		fx.Invoke(
			log.Invoke,
			apm.Invoke,
			grpc.Invoke,
			webserver.Invoke,
		),
	)
	app.LocalMessage()
	application.Run()
}
