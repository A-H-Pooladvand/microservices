package configs

import "po/pkg/vault"

type Redis struct {
	Addr string `env:"REDIS_ADDRESS" envDefault:"127.0.0.1:5672" json:"address"`
	User string `env:"REDIS_USER,notEmpty" json:"user"`
	Pass string `env:"REDIS_PASS" json:"pass"`
}

func NewRedis(client *vault.Client) (*Redis, error) {
	c := &Redis{}

	err := Parse("redis", c, client)

	return c, err
}
