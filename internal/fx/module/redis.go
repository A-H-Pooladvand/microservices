package module

import (
	"context"
	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/pkg/redis"
	"go.uber.org/fx"
)

var Redis = fx.Module("redis", fx.Provide(
	NewRedis,
))

func NewRedis(lc fx.Lifecycle, config *config.Config) *redis.Redis {
	client := redis.New(redis.Config{
		Address:  config.Redis.Addr,
		User:     config.Redis.User,
		Password: config.Redis.Pass,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}
