package cfg

import "github.com/caarlos0/env/v10"

type App struct {
	Name string `env:"APP_NAME,notEmpty" envDefault:"App"`
	Port string `env:"APP_PORT" envDefault:"8000"`
}

func NewApp() App {
	cfg := App{}

	if err := env.Parse(&cfg); err != nil {
		// Todo:: Handle Panic!
		panic(err)
	}

	return cfg
}
