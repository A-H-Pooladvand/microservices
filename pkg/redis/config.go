package redis

import (
	"time"
)

// Config holds the Redis connection configuration.
type Config struct {
	Address         string
	User            string
	Password        string
	DB              int
	MaxRetries      int
	PoolSize        int
	MinIdleConns    int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

// Option is a function that configures the Config.
type Option func(*Config)

// WithDB sets the database number.
func WithDB(db int) Option {
	return func(c *Config) {
		c.DB = db
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(retries int) Option {
	return func(c *Config) {
		c.MaxRetries = retries
	}
}

// WithPoolSize sets the pool size.
func WithPoolSize(size int) Option {
	return func(c *Config) {
		c.PoolSize = size
	}
}

// WithMinIdleConns sets the minimum number of idle connections.
func WithMinIdleConns(conns int) Option {
	return func(c *Config) {
		c.MinIdleConns = conns
	}
}

// WithDialTimeout sets the dial timeout.
func WithDialTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.DialTimeout = timeout
	}
}

// WithReadTimeout sets the read timeout.
func WithReadTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets the write timeout.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.WriteTimeout = timeout
	}
}

// NewConfig creates a new Config instance with default values.
func NewConfig(address, user, password string, opts ...Option) Config {
	cfg := Config{
		Address:         address,
		User:            user,
		Password:        password,
		DB:              0,
		MaxRetries:      3,
		PoolSize:        10,
		MinIdleConns:    3,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Address == "" {
		return ErrInvalidAddress
	}
	return nil
}
