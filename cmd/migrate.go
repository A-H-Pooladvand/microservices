package cmd

import (
	"context"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/internal/fx/invoke"
	"github.com/a-h-pooladvand/microservices/internal/fx/module"
	"github.com/a-h-pooladvand/microservices/internal/log"
	"github.com/a-h-pooladvand/microservices/internal/model"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs the migrations",
	Run:   runMigrations,
}

func runMigrations(cmd *cobra.Command, args []string) {
	application := fx.New(
		module.Vault,
		module.Postgres,

		fx.Invoke(
			invoke.Log,
			//apm.Invoke,
			func(db *gorm.DB) {
				err := db.AutoMigrate(
					model.Models...,
				)

				if err != nil {
					log.Fatal("failed to run the migrations", zap.Error(err))
				}
			},
		),
		fx.Provide(
			config.New,
		),
	)

	if err := application.Start(context.Background()); err != nil {
		log.Fatal("failed to start the application", zap.Error(err))

		return
	}
	color.Green("All migrations completed successfully")
}
