package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/fx"
	"po/configs"
	"po/internal/app"
	"po/internal/db"
	"po/internal/etcd"
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
			configs.NewEtcd,
			vault.NewConfig,
			configs.NewPostgres,
			db.New,
			// Loading services
			logstash.New,
			vault.Provide,
			rabbitmq.Provide,
			trace.Provide,
			prometheus.Provide,
			etcd.Provide,

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
			func(config *configs.App) {
				color.Green("⇨ http server started on http://127.0.0.1:%v\n", config.AppPort)
				color.Green("gRPC server started on [::]:" + config.GrpcPort)
			},
			func(c *clientv3.Client) {
				r, err := c.Get(context.TODO(), "demo")

				fmt.Println(r, err)
			},
		),
	)
	app.LocalMessage()
	application.Run()
}
