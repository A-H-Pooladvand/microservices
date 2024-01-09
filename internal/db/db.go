package db

import (
	"go.uber.org/zap"
	"po/cfg"
	"po/pkg/postgres"
)

func New() *postgres.Client {
	config := cfg.NewPostgres()

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

	if err != nil {
		zap.L().Panic("err", zap.Error(err))
	}

	return c
}
