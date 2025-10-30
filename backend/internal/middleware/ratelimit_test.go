package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	t.Cleanup(func() {
		client.Close()
		mr.Close()
	})

	return client
}

func TestRedisRateLimiter_RateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		rate           int
		window         time.Duration
		requests       int
		expectedStatus []int
		description    string
	}{
		{
			name:           "allow_requests_within_limit",
			rate:           3,
			window:         time.Minute,
			requests:       3,
			expectedStatus: []int{200, 200, 200},
			description:    "should allow requests within the rate limit",
		},
		{
			name:           "block_requests_over_limit",
			rate:           2,
			window:         time.Minute,
			requests:       4,
			expectedStatus: []int{200, 200, 429, 429},
			description:    "should block requests exceeding the rate limit",
		},
		{
			name:           "sliding_window_behavior",
			rate:           2,
			window:         100 * time.Millisecond,
			requests:       3,
			expectedStatus: []int{200, 200, 429},
			description:    "should enforce sliding window behavior",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient := setupTestRedis(t)
			limiter := NewRedisRateLimiter(redisClient)

			router := gin.New()
			router.Use(limiter.RateLimit(tt.rate, tt.window))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "ok"})
			})

			var responses []*httptest.ResponseRecorder
			for i := 0; i < tt.requests; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "127.0.0.1:12345" // Set consistent IP
				router.ServeHTTP(w, req)
				responses = append(responses, w)

				if i < tt.requests-1 && tt.window > 10*time.Millisecond {
					time.Sleep(10 * time.Millisecond) // Small delay between requests
				}
			}

			require.Len(t, responses, tt.requests, "should have correct number of responses")
			for i, resp := range responses {
				expectedStatus := tt.expectedStatus[i]
				assert.Equal(t, expectedStatus, resp.Code,
					"request %d should have status %d, got %d", i+1, expectedStatus, resp.Code)
			}
		})
	}
}

func TestRedisRateLimiter_Headers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	redisClient := setupTestRedis(t)
	limiter := NewRedisRateLimiter(redisClient)

	router := gin.New()
	router.Use(limiter.RateLimit(2, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// First request
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w1, req1)

	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, "2", w1.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "1", w1.Header().Get("X-RateLimit-Remaining"))

	// Second request
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, "2", w2.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", w2.Header().Get("X-RateLimit-Remaining"))

	// Third request (should be blocked)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w3, req3)

	assert.Equal(t, 429, w3.Code)
	assert.Equal(t, "2", w3.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", w3.Header().Get("X-RateLimit-Remaining"))
	assert.Equal(t, "60", w3.Header().Get("Retry-After"))
}

func TestRedisRateLimiter_DifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	redisClient := setupTestRedis(t)
	limiter := NewRedisRateLimiter(redisClient)

	router := gin.New()
	router.Use(limiter.RateLimit(1, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Request from IP 1
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// Request from IP 2 (should be allowed)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.2:12345"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)

	// Second request from IP 1 (should be blocked)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 429, w3.Code)
}

func TestRedisRateLimiter_DifferentEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	redisClient := setupTestRedis(t)
	limiter := NewRedisRateLimiter(redisClient)

	router := gin.New()
	router.Use(limiter.RateLimit(1, time.Minute))
	router.GET("/api1", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api1"})
	})
	router.GET("/api2", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api2"})
	})

	// Request to /api1
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api1", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// Request to /api2 (should be allowed, different endpoint)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api2", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)

	// Second request to /api1 (should be blocked)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api1", nil)
	req3.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 429, w3.Code)
}

func TestRedisRateLimiter_RedisError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a Redis client that will fail
	redisClient := redis.NewClient(&redis.Options{
		Addr: "invalid:address",
	})

	limiter := NewRedisRateLimiter(redisClient)

	router := gin.New()
	router.Use(limiter.RateLimit(1, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w, req)

	// Should continue with fallback header when Redis fails
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "true", w.Header().Get("X-RateLimit-Fallback"))
}

func TestRedisRateLimiter_WindowExpiration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	redisClient := setupTestRedis(t)
	limiter := &RedisRateLimiter{
		redis:     redisClient,
		keyPrefix: "test-ratelimit", // Use unique prefix for this test
	}

	router := gin.New()
	router.Use(limiter.RateLimit(1, 3*time.Second)) // Use 3 second window
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// First request
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// Second request (should be blocked)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 429, w2.Code)

	// Wait for window to expire (3 seconds + buffer)
	time.Sleep(3100 * time.Millisecond)

	// Small additional delay to ensure timing is reliable
	time.Sleep(100 * time.Millisecond)

	// Third request (should be allowed after window expires)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "127.0.0.1:12345"
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
}
