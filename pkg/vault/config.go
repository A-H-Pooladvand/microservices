package vault

import (
	"errors"
	"time"
)

// Config holds the Vault connection configuration.
type Config struct {
	Address        string
	RoleID         string
	SecretID       string
	MountPath      string
	Namespace      string
	Timeout        time.Duration
	MaxRetries     int
	RetryWaitMin   time.Duration
	RetryWaitMax   time.Duration
	TLSEnabled     bool
	TLSSkipVerify  bool
	CACert         string
	ClientCert     string
	ClientKey      string
}

// Option is a function that configures the Config.
type Option func(*Config)

// WithMountPath sets the mount path.
func WithMountPath(path string) Option {
	return func(c *Config) {
		c.MountPath = path
	}
}

// WithNamespace sets the namespace.
func WithNamespace(ns string) Option {
	return func(c *Config) {
		c.Namespace = ns
	}
}

// WithTimeout sets the timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(n int) Option {
	return func(c *Config) {
		c.MaxRetries = n
	}
}

// WithTLS enables TLS.
func WithTLS(enabled bool) Option {
	return func(c *Config) {
		c.TLSEnabled = enabled
	}
}

// WithTLSSkipVerify sets whether to skip TLS verification.
func WithTLSSkipVerify(skip bool) Option {
	return func(c *Config) {
		c.TLSSkipVerify = skip
	}
}

// WithCACert sets the CA certificate path.
func WithCACert(path string) Option {
	return func(c *Config) {
		c.CACert = path
	}
}

// NewConfig creates a new Config instance with default values.
func NewConfig(address, roleID, secretID string, opts ...Option) Config {
	cfg := Config{
		Address:       address,
		RoleID:        roleID,
		SecretID:      secretID,
		MountPath:     "secret",
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryWaitMin:  500 * time.Millisecond,
		RetryWaitMax:  5 * time.Second,
		TLSEnabled:    false,
		TLSSkipVerify: false,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Address == "" {
		return errors.New("vault: address is required")
	}
	return nil
}

// Empty returns true if the configuration is empty.
func (c *Config) Empty() bool {
	return c.Address == "" && c.RoleID == "" && c.SecretID == ""
}
