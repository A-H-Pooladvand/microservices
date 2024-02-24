package db

import (
	"go.uber.org/zap"
	"po/configs"
	"po/pkg/postgres"
)

func New() *postgres.Client {
	config, err := configs.NewPostgres()

	if err != nil {
		zap.L().Panic("err", zap.Error(err))
	}

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
