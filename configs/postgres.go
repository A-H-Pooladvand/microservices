package configs

import "po/pkg/vault"

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"127.0.0.1" json:"host"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432" json:"port"`
	Username string `env:"POSTGRES_USERNAME,notEmpty" json:"username"`
	Password string `env:"POSTGRES_PASSWORD" json:"password"`
	DB       string `env:"POSTGRES_DB" json:"DB"`
	Timeout  int    `env:"POSTGRES_TIMEOUT" json:"timeout"`
}

func NewPostgres(client *vault.Client) (*Postgres, error) {
	c := &Postgres{}

	err := Parse("database", c, client)

	return c, err
}
