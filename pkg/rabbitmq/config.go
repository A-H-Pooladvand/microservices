package rabbitmq

import (
	"fmt"
	"time"
)

// Config holds the RabbitMQ connection configuration.
type Config struct {
	Address          string
	User             string
	Password         string
	VHost            string
	Heartbeat        time.Duration
	ConnectionName   string
	ReconnectDelay   time.Duration
	MaxReconnects    int
	PrefetchCount    int
	PrefetchGlobal   bool
}

// Option is a function that configures the Config.
type Option func(*Config)

// WithVHost sets the virtual host.
func WithVHost(vhost string) Option {
	return func(c *Config) {
		c.VHost = vhost
	}
}

// WithHeartbeat sets the heartbeat interval.
func WithHeartbeat(d time.Duration) Option {
	return func(c *Config) {
		c.Heartbeat = d
	}
}

// WithConnectionName sets the connection name.
func WithConnectionName(name string) Option {
	return func(c *Config) {
		c.ConnectionName = name
	}
}

// WithReconnectDelay sets the reconnect delay.
func WithReconnectDelay(d time.Duration) Option {
	return func(c *Config) {
		c.ReconnectDelay = d
	}
}

// WithMaxReconnects sets the maximum number of reconnects.
func WithMaxReconnects(n int) Option {
	return func(c *Config) {
		c.MaxReconnects = n
	}
}

// WithPrefetchCount sets the prefetch count.
func WithPrefetchCount(count int) Option {
	return func(c *Config) {
		c.PrefetchCount = count
	}
}

// NewConfig creates a new Config instance with default values.
func NewConfig(address, user, password string, opts ...Option) Config {
	cfg := Config{
		Address:        address,
		User:           user,
		Password:       password,
		VHost:          "/",
		Heartbeat:      10 * time.Second,
		ConnectionName: "microservices",
		ReconnectDelay: 5 * time.Second,
		MaxReconnects:  10,
		PrefetchCount:  10,
		PrefetchGlobal: false,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// DSN returns the AMQP connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s%s",
		c.User,
		c.Password,
		c.Address,
		c.VHost,
	)
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Address == "" {
		return ErrInvalidAddress
	}
	return nil
}
