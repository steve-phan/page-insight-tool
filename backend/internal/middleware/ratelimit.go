package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements IP-based rate limiting using a simple counter approach
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
	cleanup  *time.Ticker
}

type visitor struct {
	lastSeen time.Time
	count    int
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
		// Cleanup interval: use window duration, but at least 1 minute, at most 5 minutes
		cleanup: time.NewTicker(calculateCleanupInterval(window)),
	}

	// Start cleanup goroutine to prevent memory leaks
	go rl.cleanupVisitors()

	return rl
}

// calculateCleanupInterval determines the cleanup interval based on window size
// It should be reasonable relative to the window to avoid frequent unnecessary checks
func calculateCleanupInterval(window time.Duration) time.Duration {
	// Use window * 2 as cleanup interval, but bound it between 1 minute and 5 minutes
	cleanupInterval := window * 2

	if cleanupInterval < 1*time.Minute {
		return 1 * time.Minute
	}
	if cleanupInterval > 5*time.Minute {
		return 5 * time.Minute
	}

	return cleanupInterval
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			v = &visitor{
				lastSeen: time.Now(),
				count:    0,
			}
			rl.visitors[ip] = v
		}
		rl.mu.Unlock()

		v.mu.Lock()
		defer v.mu.Unlock()

		// Reset counter if window expired
		if time.Since(v.lastSeen) > rl.window {
			v.count = 0
			v.lastSeen = time.Now()
		}

		// Check if limit exceeded
		if v.count >= rl.rate {
			v.lastSeen = time.Now()

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.rate))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", fmt.Sprintf("%d", int(rl.window.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Increment counter
		v.count++
		v.lastSeen = time.Now()

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.rate))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.rate-v.count))

		c.Next()
	}
}

// cleanupVisitors removes old visitors to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			if time.Since(v.lastSeen) > rl.window*2 {
				delete(rl.visitors, ip)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}
