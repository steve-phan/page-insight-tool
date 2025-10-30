package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// RedisRateLimiter provides rate limiting for multiple endpoints with different limits
type RedisRateLimiter struct {
	redis     *redis.Client
	keyPrefix string
}

// NewRedisRateLimiter creates a centralized rate limiter service
func NewRedisRateLimiter(redisClient *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		redis:     redisClient,
		keyPrefix: "ratelimit",
	}
}

// RateLimit creates middleware for a specific rate limit configuration
func (rls *RedisRateLimiter) RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		endpoint := c.Request.URL.Path
		key := fmt.Sprintf("%s:%s:%s", rls.keyPrefix, ip, endpoint)

		ctx := context.Background()
		now := time.Now()

		requestId := fmt.Sprintf("%d:%s", now.UnixNano(), uuid.New().String()[:8])

		// Create a Redis pipeline to execute multiple commands atomically
		pipe := rls.redis.Pipeline()

		// Add the new request to the sorted set
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  float64(now.UnixMilli()),
			Member: requestId,
		})

		// Remove entries outside the time window
		pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", float64(now.Add(-window).UnixMilli())))

		// Count the remaining requests in the window
		pipe.ZCard(ctx, key)

		// Set TTL for automatic cleanup
		pipe.Expire(ctx, key, window*2)

		cmds, err := pipe.Exec(ctx)
		if err != nil {
			fmt.Printf("Redis rate limiter error: %v\n", err)
			c.Header("X-RateLimit-Fallback", "true")
			c.Next()
			return
		}

		count := cmds[2].(*redis.IntCmd).Val()

		if count > int64(rate) {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rate))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
			c.JSON(429, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		remaining := rate - int(count)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rate))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}
