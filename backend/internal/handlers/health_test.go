package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/services/health"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		App: config.AppConfig{
			Name:        "Test App",
			Version:     "1.0.0",
			Environment: "test",
		},
	}

	// Create health service
	healthService := health.NewHealthService(cfg)

	// Create test router
	router := gin.New()
	router.GET("/health", HealthHandler(healthService, cfg))

	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	// Parse response body
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check required fields - updated to match new health service response
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "test", response["environment"])
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "version")
	assert.Contains(t, response, "build_date")
}

func TestHealthHandlerWithDifferentEnvironments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	environments := []string{"development", "staging", "production"}

	for _, env := range environments {
		t.Run("environment_"+env, func(t *testing.T) {
			cfg := &config.Config{
				App: config.AppConfig{
					Environment: env,
				},
			}

			// Create health service
			healthService := health.NewHealthService(cfg)

			router := gin.New()
			router.GET("/health", HealthHandler(healthService, cfg))

			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, env, response["environment"])
		})
	}
}

func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}
	// Create health service
	healthService := health.NewHealthService(cfg)

	router := gin.New()
	router.GET("/health", HealthHandler(healthService, cfg))

	// Test POST request (should return 404)
	req, _ := http.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
