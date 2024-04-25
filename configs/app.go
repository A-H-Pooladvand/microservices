package configs

import "po/pkg/vault"

type App struct {
	Name     string `env:"APP_NAME,notEmpty" envDefault:"App" json:"name"`
	AppPort  string `env:"APP_PORT" envDefault:"8000" json:"port"`
	GrpcPort string `env:"APP_GRPC_PORT" envDefault:"8500" json:"grpc_port"`
	Debug    string `env:"APP_DEBUG" envDefault:"true" json:"debug"`
}

func NewApp(client *vault.Client) (*App, error) {
	c := &App{}

	return c, Parse("app", c, client)
}
