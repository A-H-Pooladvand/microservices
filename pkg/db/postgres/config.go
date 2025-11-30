package postgres

import (
	"fmt"
	"time"
)

// Config holds the PostgreSQL connection configuration.
type Config struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DB              string
	Timeout         time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	SSLMode         string
	TimeZone        string
}

// Option is a function that configures the Config.
type Option func(*Config)

// WithMaxOpenConns sets the maximum number of open connections.
func WithMaxOpenConns(n int) Option {
	return func(c *Config) {
		c.MaxOpenConns = n
	}
}

// WithMaxIdleConns sets the maximum number of idle connections.
func WithMaxIdleConns(n int) Option {
	return func(c *Config) {
		c.MaxIdleConns = n
	}
}

// WithConnMaxLifetime sets the maximum connection lifetime.
func WithConnMaxLifetime(d time.Duration) Option {
	return func(c *Config) {
		c.ConnMaxLifetime = d
	}
}

// WithSSLMode sets the SSL mode.
func WithSSLMode(mode string) Option {
	return func(c *Config) {
		c.SSLMode = mode
	}
}

// WithTimeZone sets the timezone.
func WithTimeZone(tz string) Option {
	return func(c *Config) {
		c.TimeZone = tz
	}
}

// NewConfig creates a new Config instance with default values.
func NewConfig(
	host string,
	port string,
	username string,
	password string,
	db string,
	timeout int,
	opts ...Option,
) Config {
	cfg := Config{
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DB:              db,
		Timeout:         time.Second * time.Duration(timeout),
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		SSLMode:         "disable",
		TimeZone:        "Asia/Tehran",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.DB,
		c.SSLMode,
		c.TimeZone,
	)
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if c.DB == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}
