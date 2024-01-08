package db

import (
	"po/cfg"
	"po/pkg/postgres"
	"po/pkg/zlog"
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
		zlog.Panic(err)
	}

	return c
}
