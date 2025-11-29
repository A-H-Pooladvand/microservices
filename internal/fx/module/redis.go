package module

import (
	"context"
	"fmt"

	"github.com/a-h-pooladvand/microservices/config"
	"github.com/a-h-pooladvand/microservices/pkg/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Redis = fx.Module("redis", fx.Provide(
	NewRedis,
))

// NewRedis creates a new Redis client with observability.
func NewRedis(lc fx.Lifecycle, cfg *config.Config) (*redis.Client, error) {
	client, err := redis.New(
		redis.NewConfig(
			cfg.Redis.Addr,
			cfg.Redis.User,
			cfg.Redis.Pass,
		),
		zap.L(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
