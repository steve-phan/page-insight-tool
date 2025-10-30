package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/steve-phan/page-insight-tool/internal/config"
)

// RedisService provides Redis client functionality
type RedisService struct {
	client *redis.Client
	config *config.Config
}

// NewRedisService creates a new Redis service
func NewRedisService(cfg *config.Config) (*RedisService, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisService{
		client: rdb,
		config: cfg,
	}, nil
}

// NewTestRedisService creates a Redis service for testing (no actual connection)
func NewTestRedisService(cfg *config.Config) *RedisService {
	// For testing, create a client that won't actually connect
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   cfg.Redis.DB,
	})

	return &RedisService{
		client: rdb,
		config: cfg,
	}
}

// GetClient returns the Redis client
func (r *RedisService) GetClient() *redis.Client {
	return r.client
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}

// IsConnected checks if Redis is connected
func (r *RedisService) IsConnected(ctx context.Context) bool {
	return r.client.Ping(ctx).Err() == nil
}
