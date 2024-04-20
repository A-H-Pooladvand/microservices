package redis

import (
	"go.uber.org/zap"
	"po/configs"
)

type Config struct {
	Address  string
	User     string
	Password string
}

func NewDefaultConfig() Config {
	c, err := configs.NewRedis()

	if err != nil {
		zap.L().Fatal("redis connection fail", zap.Error(err))
	}

	return Config{
		Address:  c.Addr,
		User:     c.User,
		Password: c.Pass,
	}
}
