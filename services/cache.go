package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService handles Redis caching operations
type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new instance of CacheService
func NewCacheService(redisURL string) *CacheService {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// Fallback to default localhost connection
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	client := redis.NewClient(opt)
	return &CacheService{
		client: client,
	}
}

// Set stores a key-value pair in Redis with TTL
func (c *CacheService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value from Redis by key
func (c *CacheService) Get(ctx context.Context, key string) (string, error) {
	result := c.client.Get(ctx, key)
	if result.Err() == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	return result.Val(), result.Err()
}

// Delete removes a key from Redis
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (c *CacheService) Close() error {
	return c.client.Close()
}

// Ping tests the Redis connection
func (c *CacheService) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
