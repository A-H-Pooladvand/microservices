package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"po/configs"
)

type Client struct {
	*amqp.Connection
}

func New(c Config) (*Client, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		c.User,
		c.Password,
		c.Address,
	)

	conn, err := amqp.Dial(dsn)

	if err != nil {
		return nil, err
	}

	return &Client{
		Connection: conn,
	}, nil
}

func Provide(lc fx.Lifecycle, c *configs.RabbitMQ) *Client {
	conn, err := New(Config{
		Address:  c.Addr,
		User:     c.User,
		Password: c.Pass,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return err
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn
}
