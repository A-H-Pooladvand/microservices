package configs

import "po/pkg/vault"

type RabbitMQ struct {
	Addr string `env:"RABBITMQ_ADDRESS" envDefault:"127.0.0.1:15672" json:"address"`
	User string `env:"RABBITMQ_USER,notEmpty" json:"user"`
	Pass string `env:"RABBITMQ_PASS" json:"pass"`
}

func NewRabbitMQ(client *vault.Client) (*RabbitMQ, error) {
	c := &RabbitMQ{}

	err := Parse("rabbitmq", c, client)

	return c, err
}
