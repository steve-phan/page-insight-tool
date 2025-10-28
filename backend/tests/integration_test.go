package tests

import (
	"net/http"
	"testing"
	"time"

	"page-insight-tool/internal/config"
	"page-insight-tool/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const configPath = "../config/config.yaml"

func TestServerLifecycle(t *testing.T) {
	// Load test configuration
	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	// Create server
	srv := server.New(cfg)
	require.NotNil(t, srv)

	// Test server start
	err = srv.Start()
	assert.NoError(t, err)

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	t.Run("health_endpoint", func(t *testing.T) {
		resp, err := http.Get("http://" + cfg.GetAddress() + "/api/v1/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
	})

	// Test server stop
	err = srv.Stop()
	assert.NoError(t, err)
}

func TestServerGracefulShutdown(t *testing.T) {
	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	srv := server.New(cfg)
	err = srv.Start()
	require.NoError(t, err)

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test that server responds
	resp, err := http.Get("http://" + cfg.GetAddress() + "/api/v1/health")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Stop server
	err = srv.Stop()
	assert.NoError(t, err)

	// Test that server no longer responds
	time.Sleep(100 * time.Millisecond)
	_, err = http.Get("http://" + cfg.GetAddress() + "/api/v1/health")
	assert.Error(t, err) // Should fail because server is stopped
}

func TestServerWithInvalidConfig(t *testing.T) {
	// Test with invalid port
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 0, // Invalid port
		},
	}

	// This should not panic, but server might fail to start
	srv := server.New(cfg)
	assert.NotNil(t, srv)
}

func TestConcurrentRequests(t *testing.T) {
	cfg, err := config.LoadConfig(configPath)
	require.NoError(t, err)

	srv := server.New(cfg)
	err = srv.Start()
	require.NoError(t, err)
	defer srv.Stop()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Make concurrent requests
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			resp, err := http.Get("http://" + cfg.GetAddress() + "/api/v1/health")
			if err == nil {
				resp.Body.Close()
			}
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
