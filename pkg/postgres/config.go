package postgres

import "time"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       string
	Timeout  time.Duration
}

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
