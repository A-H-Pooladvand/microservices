package cmd

import (
	"context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"po/configs"
	"po/internal/apm"
	"po/internal/app"
	"po/internal/db"
	"po/internal/models"
	"po/internal/vault"
	"po/pkg/log"
	"po/pkg/logstash"
	"po/pkg/postgres"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs the migrations",
	Run:   runMigrations,
}

func runMigrations(cmd *cobra.Command, args []string) {
	app.LoadEnvironmentVariablesInLocalEnv()

	application := fx.New(
		fx.Provide(
			// Loading configs
			configs.NewApp,
			configs.NewLogstash,
			configs.NewPostgres,
			vault.NewConfig,
			// Loading services
			logstash.New,
			vault.Provide,
			db.New,
		),

		fx.Invoke(
			log.Invoke,
			apm.Invoke,
			func(db *postgres.Client) {
				err := db.AutoMigrate(
					migrations()...,
				)

				if err != nil {
					zap.L().Fatal("failed to run the migrations", zap.Error(err))
				}
			},
		),
	)
	app.LocalMessage()

	if err := application.Start(context.Background()); err != nil {
		zap.L().Fatal("failed to start the application", zap.Error(err))

		return
	}
	color.Green("All migrations completed successfully")
}

func migrations() []any {
	return []any{
		models.User{},
	}
}
