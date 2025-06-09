package module

import (
	"context"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/pkg/db/postgres"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Postgres = fx.Module("postgres", fx.Provide(
	NewPostgres,
))

// NewPostgres New creates a new database connection.
func NewPostgres(lc fx.Lifecycle, config *config.Config) *gorm.DB {
	db, err := postgres.New(
		postgres.NewConfig(
			config.Postgres.Host,
			config.Postgres.Port,
			config.Postgres.Username,
			config.Postgres.Password,
			config.Postgres.DB,
			config.Postgres.Timeout,
			//config.App.Debuggable(),
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
