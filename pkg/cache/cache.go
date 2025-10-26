package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache interface for different cache implementations
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache instance
// prefix is optional - if empty, no prefix is used
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a value from cache
func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, r.prefix+key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	return val, err
}

// Set stores a value in cache with TTL
func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, r.prefix+key, value, ttl).Err()
}

// Delete removes one or more keys from cache
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.prefix + key
	}

	return r.client.Del(ctx, prefixedKeys...).Err()
}

// DeletePattern removes all keys matching a pattern
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	// Scan for all keys matching the pattern
	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = r.client.Scan(ctx, cursor, r.prefix+pattern, 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	// Delete all found keys
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in cache
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, r.prefix+key).Result()
	return count > 0, err
}

// NoOpCache is a cache implementation that does nothing (for testing or disabling cache)
type NoOpCache struct{}

func (n *NoOpCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (n *NoOpCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func (n *NoOpCache) Delete(ctx context.Context, keys ...string) error {
	return nil
}

func (n *NoOpCache) DeletePattern(ctx context.Context, pattern string) error {
	return nil
}

func (n *NoOpCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}
