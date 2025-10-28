package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"page-insight-tool/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		App: config.AppConfig{
			Name:        "Test App",
			Environment: "test",
		},
	}

	router := SetupRoutes(cfg)

	// Test health endpoint
	t.Run("health_endpoint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	// Test OPTIONS handling
	t.Run("options_request", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/api/v1/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	// Test 404 for unknown routes
	t.Run("unknown_route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/unknown", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMiddlewareSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}
	router := SetupRoutes(cfg)

	// Test that middleware is applied by checking if recovery works
	t.Run("panic_recovery", func(t *testing.T) {
		// Add a route that panics
		router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		req, _ := http.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()

		// This should not crash the test
		router.ServeHTTP(w, req)

		// Recovery middleware should handle the panic
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAPIRoutesGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}
	router := SetupRoutes(cfg)

	// Test that API routes are properly grouped under /api/v1
	apiRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/health"},
	}

	for _, route := range apiRoutes {
		t.Run(route.method+"_"+route.path, func(t *testing.T) {
			req, _ := http.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
