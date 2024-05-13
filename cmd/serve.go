package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"po/configs"
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
	"po/pkg/trace"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the application",
	// Provide the necessary protocols such as gRPC, HTTP, etc...
	Run: runApplication,
}

func runApplication(cmd *cobra.Command, args []string) {
	app.LoadEnvironmentVariablesInLocalEnv()

	application := fx.New(
		user.Module,
		fx.Provide(
			// Loading configs
			configs.NewApp,
			configs.NewLogstash,
			configs.NewRabbitMQ,
			configs.NewRedis,
			configs.NewJaeger,
			vault.NewConfig,
			configs.NewPostgres,
			db.New,
			// Loading services
			logstash.New,
			vault.Provide,
			rabbitmq.Provide,
			trace.Provide,

			handlers.NewRestHandlers,
			handlers.NewGrpcHandlers,
			fx.Annotate(
				redis.Provide,
				fx.As(new(cache.Cache)),
			),
		),

		fx.Invoke(
			log.Invoke,
			//apm.Provide,
			grpc.Invoke,
			webserver.Invoke,
		),
	)
	app.LocalMessage()
	application.Run()
}
