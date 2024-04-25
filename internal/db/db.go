package db

import (
	"context"
	"go.uber.org/fx"
	"po/configs"
	"po/pkg/postgres"
)

func New(lc fx.Lifecycle, config *configs.Postgres) *postgres.Client {
	c, err := postgres.New(
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
			return c.Close()
		},
	})

	return c
}
