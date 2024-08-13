package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"po/configs"
	"po/pkg/cache"
	"time"
)

// Verify interface compliance
var _ cache.Cache = (*Redis)(nil)

type Redis struct {
	client *redis.Client
}

func New(c Config) *Redis {
	client := redis.NewClient(&redis.Options{
		//ClientName:            "",
		//OnConnect:             nil,
		Addr:     c.Address,
		Username: c.User,
		Password: c.Password,
		DB:       0,
	})

	return &Redis{
		client: client,
	}
}

func Provide(lc fx.Lifecycle, c *configs.Redis) *Redis {
	client := New(Config{
		Address:  c.Addr,
		User:     c.User,
		Password: c.Pass,
	})

	return client
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(
		ctx,
		key,
		r.normalize(value),
		ttl,
	).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) Remember(ctx context.Context, key string, value any, ttl time.Duration) ([]byte, error) {
	if err := r.client.SetNX(
		ctx,
		key,
		r.normalize(value),
		ttl,
	).Err(); err != nil {
		return nil, err
	}

	return r.Get(ctx, key)
}

func (r *Redis) Forever(ctx context.Context, key string, value any) ([]byte, error) {
	return r.Remember(ctx, key, value, 0)
}

func (r *Redis) normalize(value any) []byte {
	switch v := value.(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		b, _ := json.Marshal(v)
		return b
	}
}
