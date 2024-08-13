package db

import (
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"po/configs"
	"po/pkg/db/postgres"
)

// New creates a new database connection.
func New(lc fx.Lifecycle, config *configs.Postgres) *gorm.DB {
	db, err := postgres.New(
		postgres.NewConfig(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
			config.DB,
			config.Timeout,
		),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return err
		},
		OnStop: func(ctx context.Context) error {
			sql, err := db.DB()

			if err != nil {
				return err
			}

			return sql.Close()
		},
	})

	return db
}
