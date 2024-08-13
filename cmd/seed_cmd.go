package cmd

import (
	"context"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"po/cmd/seed"
	"po/configs"
	"po/internal/app"
	"po/internal/db"
	"po/internal/vault"
	"po/pkg/log"
	"po/pkg/logstash"
)

var seeders = []seed.Seeder{
	seed.UserSeeder{},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seeds the database",
	Run:   runSeeders,
}

func runSeeders(cmd *cobra.Command, args []string) {
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
			//apm.Invoke,
			func(db *gorm.DB) {
				for _, seeder := range seeders {
					seeder.Run(db)
				}
			},
		),
	)
	app.LocalMessage()

	if err := application.Start(context.Background()); err != nil {
		zap.L().Fatal("failed to seed the database", zap.Error(err))

		return
	}
	color.Green("All seeders completed successfully")
}
