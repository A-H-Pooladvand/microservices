package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type Redis struct {
	*redis.Client
	lock *redislock.Client
}

func New(c Config) *Redis {
	client := redis.NewClient(&redis.Options{
		//ClientName:            "",
		//OnConnect:             nil,
		Addr:     c.Address,
		Username: c.User,
		Password: c.Password,
		//DB:       0,
	})

	r := &Redis{
		Client: client,
		lock:   redislock.New(client),
	}

	return r
}

// Transaction acquires a distributed lock for the given key with retry attempts,
// executes the provided function f, and then releases the lock.
// It handles lock acquisition failures with retries and ensures the lock is
// released using defer.
//
// key: The key to lock.
// ttl: The maximum time the lock should be held.
// opt: Optional settings for acquiring the lock (applied to each retry attempt).
// f: The function to execute while holding the lock. It receives
//
//	the Redis client instance and should return an error if
//	the operation within the transaction fails.
//
// Returns redislock.ErrNotObtained if the lock could not be acquired
// after all retry attempts, or any other error encountered during
// lock acquisition or the execution of the provided function.
func (r *Redis) Transaction(
	ctx context.Context,
	key string,
	f func(tx *Redis) error,
) error {
	ctx, cancel := context.WithDeadline(
		ctx,
		time.Now().Add(time.Minute*2),
	)

	defer cancel()
	lock, err := r.lock.Obtain(ctx, key+"-lock", time.Minute*2, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(500*time.Millisecond), 10),
	})

	if err != nil {
		return err
	}

	// If we are here, the lock was successfully obtained (err == nil)

	// Ensure the lock is released when the function exits.
	// We capture the error from Release to handle it separately.
	defer func() {
		releaseErr := lock.Release(ctx)
		if releaseErr != nil {
			// Log a warning or error if releasing the lock fails.
			// This is important because the critical section might
			// have succeeded, but the cleanup failed.
			log.Printf("failed to release lock for key %s: %v\n", key, releaseErr)
		}
	}()

	// Return the error from the executed function
	return f(r)
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) *redis.StatusCmd {
	normalizedValue, err := r.normalize(value)

	if err != nil {
		return redis.NewStatusResult("", err)
	}

	return r.Client.Set(
		ctx,
		key,
		normalizedValue,
		ttl,
	)
}

func (r *Redis) Exist(ctx context.Context, key string) (bool, error) {
	cmd := r.Client.Exists(ctx, key)

	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	return cmd.Val() == 1, nil
}

func (r *Redis) Remember(ctx context.Context, key string, ttl time.Duration, f func() (any, error)) *redis.StringCmd {
	if cmd := r.Get(ctx, key); cmd.Err() == nil {
		return cmd
	} else if !errors.Is(cmd.Err(), redis.Nil) {
		return cmd // Return any error that's not a cache miss
	}

	v, err := f()

	if err != nil {
		return redis.NewStringResult("", err)
	}

	normalizedValue, err := r.normalize(v)

	if err != nil {
		return redis.NewStringResult("", err)
	}

	result := r.Set(ctx, key, normalizedValue, ttl)

	if result.Err() != nil {
		return redis.NewStringResult("", result.Err())
	}

	// Convert the value to string and return it directly
	// This avoids an unnecessary round trip to Redis
	return redis.NewStringResult(string(normalizedValue), nil)
}

func (r *Redis) Forever(ctx context.Context, key string, f func() (any, error)) *redis.StringCmd {
	return r.Remember(ctx, key, 0, f)
}

func (r *Redis) normalize(value any) ([]byte, error) {
	// 1. Handle nil explicitly
	if value == nil {
		return nil, nil // Represent nil as nil []byte (Redis command might store "(nil)" or empty)
		// Alternatively: return []byte{}, nil // Represent nil as empty byte slice
	}

	// 2. Check for standard marshaling interfaces first for custom type handling
	// If a type implements these, it knows best how to represent itself.
	if marshaler, ok := value.(encoding.BinaryMarshaler); ok {
		data, err := marshaler.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("normalize: failed marshaling with encoding.BinaryMarshaler: %w", err)
		}
		return data, nil
	}
	// time.Time implements TextMarshaler, good example
	if marshaler, ok := value.(encoding.TextMarshaler); ok {
		data, err := marshaler.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("normalize: failed marshaling with encoding.TextMarshaler: %w", err)
		}
		return data, nil
	}

	// 3. Handle common types directly for performance and Redis compatibility
	switch v := value.(type) {
	case []byte:
		// If the input slice might be modified later by the caller,
		// you might want to return a copy:
		// clone := make([]byte, len(v))
		// copy(clone, v)
		// return clone, nil
		// For max performance (assuming caller won't modify):
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		// AppendInt is efficient; appends to nil slice, avoids intermediate string
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
	case uint8: // byte is an alias for uint8
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float64:
		// 'g' format is usually suitable, -1 precision means shortest necessary
		return strconv.AppendFloat(nil, v, 'g', -1, 64), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'g', -1, 32), nil
	case bool:
		if v {
			return []byte("1"), nil // Common representation for true
		}
		return []byte("0"), nil // Common representation for false
	// --- Optional: Handle pointers to basic types ---
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
	// ... add other pointer types (*int, *float64, *bool etc.) if needed

	// 4. Fallback to JSON for complex types (structs, maps, slices/arrays, etc.)
	default:
		// fmt.Printf("Normalizing type %T using JSON\n", v) // Debugging line
		b, err := json.Marshal(v)
		if err != nil {
			// IMPORTANT: Return the error instead of ignoring it!
			return nil, fmt.Errorf("normalize: failed marshaling to JSON for type %T: %w", v, err)
		}
		return b, nil
	}
}
