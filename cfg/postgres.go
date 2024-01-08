package cfg

import (
	"github.com/caarlos0/env/v10"
)

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"127.0.0.1"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	Username string `env:"POSTGRES_USERNAME,notEmpty"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
	Timeout  int    `env:"POSTGRES_TIMEOUT"`
}

func NewPostgres() Postgres {
	cfg := Postgres{}

	if err := env.Parse(&cfg); err != nil {
		// Todo:: Handle Panic!
		panic(err)
	}

	return cfg
}
