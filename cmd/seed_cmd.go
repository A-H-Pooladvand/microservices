package cmd

import (
	"context"
	"github.com/a-h-pooladvand/microservices/cmd/seed"
	"github.com/a-h-pooladvand/microservices/internal/fx/invoke"
	"github.com/a-h-pooladvand/microservices/internal/fx/module"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var seeders []seed.Seeder

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seeds the database",
	Run:   runSeeders,
}

func runSeeders(cmd *cobra.Command, args []string) {
	application := fx.New(
		module.Vault,
		module.Postgres,

		fx.Invoke(
			invoke.Log,
			//apm.Invoke,
			func(db *gorm.DB) {
				for _, seeder := range seeders {
					seeder.Run(db)
				}
			},
		),
	)

	if err := application.Start(context.Background()); err != nil {
		log.Fatal("failed to seed the database", zap.Error(err))

		return
	}
	color.Green("All seeders completed successfully")
}
