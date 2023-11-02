package configs

import (
	"github.com/caarlos0/env/v10"
	"po/pkg/logger"
)

type App struct {
	Port  string `env:"APP_PORT" envDefault:"8000"`
	Debug string `env:"APP_DEBUG"`
}

func NewApp() App {
	cfg := App{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error(err)
		panic("unable to load mysql env file" + err.Error())
	}

	return cfg
}
