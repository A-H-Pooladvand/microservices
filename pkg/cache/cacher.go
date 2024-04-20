package cache

import (
	"context"
	"time"
)

type Cache interface {
	// Set stores a key-value pair in the cache with an optional Time To Live (TTL).
	//  - ctx: Context for the operation (optional).
	//  - key: Unique key to identify the cached value.
	//  - value: The data to be cached.
	//  - ttl: The duration for which the key-value pair should be considered valid.
	//        If zero, the cache implementation's default TTL will be used.
	// Returns an error if storing the data fails.
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Get retrieves a value from the cache using its key.
	//  - ctx: Context for the operation (optional).
	//  - key: The key of the data to be retrieved.
	// Returns the retrieved value as a byte slice and an error if retrieval fails.
	Get(ctx context.Context, key string) ([]byte, error)

	// Remember retrieves a value from the cache using its key. If the key doesn't exist,
	// it stores the provided value with the specified TTL and returns it.
	// This method is functionally equivalent to calling Get followed by Set if the key is missing.
	//  - ctx: Context for the operation (optional).
	//  - key: The key of the data to be retrieved or stored.
	//  - value: The data to be stored if the key is missing.
	//  - ttl: The duration for which the key-value pair should be considered valid.
	//        If zero, the cache implementation's default TTL will be used.
	// Returns the retrieved or stored value as a byte slice and an error if any operation fails.
	Remember(ctx context.Context, key string, value any, ttl time.Duration) ([]byte, error)

	// Forever stores a key-value pair in the cache with no expiration.
	//  - ctx: Context for the operation (optional).
	//  - key: Unique key to identify the cached value.
	//  - value: The data to be cached.
	// This method stores the data indefinitely, but be aware that:
	//  - The cache implementation may still evict the data due to memory pressure.
	//  - Manually removing the data using Delete is necessary.
	// Returns an error if storing the data fails.
	Forever(ctx context.Context, key string, value any) ([]byte, error)

	// Delete removes one or more keys from the cache.
	//  - ctx: Context for the operation (optional).
	//  - keys: The keys to be deleted from the cache.
	// Returns an error if deleting any of the keys fails.
	Delete(ctx context.Context, keys ...string) error
}
