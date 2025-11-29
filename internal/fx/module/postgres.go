package module

import (
	"context"
	"fmt"

	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/pkg/db/postgres"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Postgres = fx.Module("postgres", fx.Provide(
	NewPostgres,
	ProvideGormDB,
))

// NewPostgres creates a new database connection with observability.
func NewPostgres(lc fx.Lifecycle, cfg *config.Config) (*postgres.Client, error) {
	client, err := postgres.New(
		postgres.NewConfig(
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.DB,
			cfg.Postgres.Timeout,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

// ProvideGormDB extracts the underlying gorm.DB from the postgres client.
// This maintains backward compatibility with existing code.
func ProvideGormDB(client *postgres.Client) *gorm.DB {
	return client.DB
}
