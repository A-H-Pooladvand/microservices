package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func New(c Config) (*amqp.Connection, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		c.User,
		c.Password,
		c.Address,
	)

	return amqp.Dial(dsn)
}
