package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       string
	Timeout  time.Duration
}

// NewConfig creates a new Config instance.
func NewConfig(
	host string,
	port string,
	username string,
	password string,
	db string,
	timeout int,
) Config {
	return Config{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		DB:       db,
		Timeout:  time.Second * time.Duration(timeout),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable TimeZone=Asia/Tehran",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.DB,
	)
}
