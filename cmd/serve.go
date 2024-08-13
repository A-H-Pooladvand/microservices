package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"po/configs"
	"po/internal/app"
	"po/internal/db"
	"po/internal/grpc"
	"po/internal/handlers"
	"po/internal/handlers/metric"
	"po/internal/handlers/user"
	"po/internal/vault"
	"po/internal/webserver"
	"po/pkg/cache"
	"po/pkg/log"
	"po/pkg/logstash"
	"po/pkg/prometheus"
	"po/pkg/rabbitmq"
	"po/pkg/redis"
	"po/pkg/trace"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the application",
	Long:  "Serve the application and start the server",
	Run:   runApplication,
}

func runApplication(cmd *cobra.Command, args []string) {
	app.LoadEnvironmentVariablesInLocalEnv()

	application := fx.New(
		user.Module,
		metric.Module,
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
			prometheus.Provide,

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
