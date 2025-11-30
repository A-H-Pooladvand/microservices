package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/a-h-pooladvand/microservices/pkg/observability"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// Client wraps go-redis client with observability support.
type Client struct {
	*redis.Client
	lock    *redislock.Client
	tracer  *observability.Tracer
	metrics *observability.CacheMetrics
	logger  *zap.Logger
}

// New creates a new Redis client with observability.
func New(cfg Config, logger *zap.Logger) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Address,
		Username:        cfg.User,
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      cfg.MaxRetries,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	})

	// Initialize observability
	tracer := observability.NewTracer("redis")
	meter := observability.NewMeter("redis")
	metrics, err := observability.NewCacheMetrics(meter, "redis")
	if err != nil {
		return nil, fmt.Errorf("failed to create cache metrics: %w", err)
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	return &Client{
		Client:  client,
		lock:    redislock.New(client),
		tracer:  tracer,
		metrics: metrics,
		logger:  logger,
	}, nil
}

// Get retrieves a value by key with tracing.
func (c *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.GET",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "GET"),
		attribute.String("db.redis.key", key),
	)
	defer span.End()

	startTime := time.Now()
	cmd := c.Client.Get(ctx, key)

	c.recordMetrics(ctx, "GET", startTime, cmd.Err())

	if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
		observability.SetSpanError(span, cmd.Err())
		c.logger.Error("redis GET failed", zap.String("key", key), zap.Error(cmd.Err()))
	} else {
		observability.SetSpanOK(span)
		if errors.Is(cmd.Err(), redis.Nil) {
			c.metrics.MissesTotal.Add(ctx, 1)
		} else {
			c.metrics.HitsTotal.Add(ctx, 1)
		}
	}

	return cmd
}

// Set stores a value with tracing.
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.SET",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "SET"),
		attribute.String("db.redis.key", key),
		attribute.Int64("db.redis.ttl_ms", ttl.Milliseconds()),
	)
	defer span.End()

	startTime := time.Now()
	normalizedValue, err := c.normalize(value)
	if err != nil {
		observability.SetSpanError(span, err)
		return redis.NewStatusResult("", err)
	}

	cmd := c.Client.Set(ctx, key, normalizedValue, ttl)

	c.recordMetrics(ctx, "SET", startTime, cmd.Err())

	if cmd.Err() != nil {
		observability.SetSpanError(span, cmd.Err())
		c.logger.Error("redis SET failed", zap.String("key", key), zap.Error(cmd.Err()))
	} else {
		observability.SetSpanOK(span)
	}

	return cmd
}

// Del deletes keys with tracing.
func (c *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.DEL",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "DEL"),
		attribute.Int("db.redis.key_count", len(keys)),
	)
	defer span.End()

	startTime := time.Now()
	cmd := c.Client.Del(ctx, keys...)

	c.recordMetrics(ctx, "DEL", startTime, cmd.Err())

	if cmd.Err() != nil {
		observability.SetSpanError(span, cmd.Err())
	} else {
		observability.SetSpanOK(span)
	}

	return cmd
}

// Exist checks if a key exists with tracing.
func (c *Client) Exist(ctx context.Context, key string) (bool, error) {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.EXISTS",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "EXISTS"),
		attribute.String("db.redis.key", key),
	)
	defer span.End()

	startTime := time.Now()
	cmd := c.Client.Exists(ctx, key)

	c.recordMetrics(ctx, "EXISTS", startTime, cmd.Err())

	if cmd.Err() != nil {
		observability.SetSpanError(span, cmd.Err())
		return false, cmd.Err()
	}

	observability.SetSpanOK(span)
	return cmd.Val() == 1, nil
}

// Remember implements cache-aside pattern with tracing.
func (c *Client) Remember(ctx context.Context, key string, ttl time.Duration, f func() (any, error)) *redis.StringCmd {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.Remember",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "REMEMBER"),
		attribute.String("db.redis.key", key),
	)
	defer span.End()

	// Try to get from cache first
	if cmd := c.Get(ctx, key); cmd.Err() == nil {
		span.SetAttributes(attribute.Bool("cache_hit", true))
		return cmd
	} else if !errors.Is(cmd.Err(), redis.Nil) {
		observability.SetSpanError(span, cmd.Err())
		return cmd
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))

	// Cache miss - execute function
	v, err := f()
	if err != nil {
		observability.SetSpanError(span, err)
		return redis.NewStringResult("", err)
	}

	normalizedValue, err := c.normalize(v)
	if err != nil {
		observability.SetSpanError(span, err)
		return redis.NewStringResult("", err)
	}

	// Store in cache
	result := c.Set(ctx, key, normalizedValue, ttl)
	if result.Err() != nil {
		observability.SetSpanError(span, result.Err())
		return redis.NewStringResult("", result.Err())
	}

	observability.SetSpanOK(span)
	return redis.NewStringResult(string(normalizedValue), nil)
}

// Forever stores a value without expiration.
func (c *Client) Forever(ctx context.Context, key string, f func() (any, error)) *redis.StringCmd {
	return c.Remember(ctx, key, 0, f)
}

// Transaction acquires a distributed lock and executes function.
func (c *Client) Transaction(ctx context.Context, key string, f func(tx *Client) error) error {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.Transaction",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "LOCK"),
		attribute.String("db.redis.lock_key", key),
	)
	defer span.End()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute*2))
	defer cancel()

	lock, err := c.lock.Obtain(ctx, key+"-lock", time.Minute*2, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(500*time.Millisecond), 10),
	})
	if err != nil {
		observability.SetSpanError(span, err)
		c.logger.Error("failed to obtain lock", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("obtain lock: %w", err)
	}

	defer func() {
		if releaseErr := lock.Release(ctx); releaseErr != nil {
			c.logger.Warn("failed to release lock", zap.String("key", key), zap.Error(releaseErr))
		}
	}()

	if err := f(c); err != nil {
		observability.SetSpanError(span, err)
		return err
	}

	observability.SetSpanOK(span)
	return nil
}

// Ping checks Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	ctx, span := c.tracer.StartWithAttributes(ctx, "redis.PING",
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", "PING"),
	)
	defer span.End()

	cmd := c.Client.Ping(ctx)
	if cmd.Err() != nil {
		observability.SetSpanError(span, cmd.Err())
		return cmd.Err()
	}

	observability.SetSpanOK(span)
	return nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.Client.Close()
}

// PoolStats returns connection pool statistics.
func (c *Client) PoolStats() *redis.PoolStats {
	return c.Client.PoolStats()
}

func (c *Client) recordMetrics(ctx context.Context, operation string, startTime time.Time, err error) {
	duration := time.Since(startTime).Seconds()
	attrs := metric.WithAttributes(attribute.String("operation", operation))

	c.metrics.OperationTime.Record(ctx, duration, attrs)

	if err != nil && !errors.Is(err, redis.Nil) {
		c.metrics.ErrorsTotal.Add(ctx, 1, attrs)
	}
}

func (c *Client) normalize(value any) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	if marshaler, ok := value.(encoding.BinaryMarshaler); ok {
		data, err := marshaler.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("normalize: failed marshaling with encoding.BinaryMarshaler: %w", err)
		}
		return data, nil
	}

	if marshaler, ok := value.(encoding.TextMarshaler); ok {
		data, err := marshaler.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("normalize: failed marshaling with encoding.TextMarshaler: %w", err)
		}
		return data, nil
	}

	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, v, 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, v, 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float64:
		return strconv.AppendFloat(nil, v, 'g', -1, 64), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'g', -1, 32), nil
	case bool:
		if v {
			return []byte("1"), nil
		}
		return []byte("0"), nil
	case *string:
		if v == nil {
			return nil, nil
		}
		return []byte(*v), nil
	case *int64:
		if v == nil {
			return nil, nil
		}
		return strconv.AppendInt(nil, *v, 10), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("normalize: failed marshaling to JSON for type %T: %w", v, err)
		}
		return b, nil
	}
}
