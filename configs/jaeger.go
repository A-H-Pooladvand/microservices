package configs

import "po/pkg/vault"

type Jaeger struct {
	Addr string `env:"JAEGER_ADDR" envDefault:"127.0.0.1:4317" json:"addr"`
}

func NewJaeger(client *vault.Client) (*Jaeger, error) {
	c := &Jaeger{}

	err := Parse("jaeger", c, client)

	return c, err
}
