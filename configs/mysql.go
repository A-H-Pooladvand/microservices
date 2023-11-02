package configs

import (
	"github.com/caarlos0/env/v10"
	"po/pkg/logger"
)

type Mysql struct {
	Host     string `env:"MYSQL_HOST"`
	Port     string `env:"MYSQL_PORT"`
	Username string `env:"MYSQL_USERNAME"`
	Password string `env:"MYSQL_PASSWORD"`
	Database string `env:"MYSQL_DATABASE"`
}

func NewMysql() Mysql {
	cfg := Mysql{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error(err)
		panic("unable to load mysql env file" + err.Error())
	}

	return cfg
}
