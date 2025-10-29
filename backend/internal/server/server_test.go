package server

import (
	"os"
	"testing"
	"time"

	"page-insight-tool/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 0, // Use port 0 for testing
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

	srv, err := NewForTesting(cfg)

	require.NoError(t, err)
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.services)
	assert.NotNil(t, srv.handlers)
	assert.NotNil(t, srv.router)
	assert.NotNil(t, srv.httpSrv)
}

func TestNewProductionMode(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Environment: "production",
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	assert.NotNil(t, srv)

}

func TestGetVersion(t *testing.T) {
	// Test with environment variables
	os.Setenv("VERSION", "1.2.3")
	defer os.Unsetenv("VERSION")

	version := getVersion()
	assert.Equal(t, "1.2.3", version)

	// Test without environment variables
	os.Unsetenv("VERSION")
	version = getVersion()
	assert.Equal(t, "dev", version)
}

func TestGetBuildDate(t *testing.T) {
	// Test with environment variables
	os.Setenv("BUILD_DATE", "2023-01-01T00:00:00Z")
	defer os.Unsetenv("BUILD_DATE")

	buildDate := getBuildDate()
	assert.Equal(t, "2023-01-01T00:00:00Z", buildDate)

	// Test without environment variables - should return current time in RFC1123 format
	os.Unsetenv("BUILD_DATE")
	buildDate = getBuildDate()
	expectedBuildDate := time.Now().Local().Format(time.RFC1123)
	// Should match RFC1123 format like "Wed, 29 Oct 2025 07:50:52 CET"
	assert.Equal(t, expectedBuildDate, buildDate)

}

func TestServerStartStop(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 0, // Use port 0 to let OS assign available port
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// Test start
	err = srv.Start()
	assert.NoError(t, err)

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test stop
	err = srv.Stop()
	assert.NoError(t, err)
}

func TestServerStopWithoutStart(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 0,
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// Test stop without start (should not panic)
	err = srv.Stop()
	assert.NoError(t, err)
}

func TestServerMultipleStops(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 0,
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	err = srv.Start()
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Stop multiple times (should not panic)
	err = srv.Stop()
	assert.NoError(t, err)

	err = srv.Stop()
	assert.NoError(t, err)
}

func TestServerWithTimeout(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         0,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	err = srv.Start()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = srv.Stop()
	assert.NoError(t, err)
}

func TestServerAddress(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// Test that the server address is correctly set
	expectedAddr := "127.0.0.1:8080"
	assert.Equal(t, expectedAddr, srv.httpSrv.Addr)
}

func TestServerHandler(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Environment: "test",
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// Test that the handler is set
	assert.NotNil(t, srv.httpSrv.Handler)
}

func TestServerTimeouts(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}

	srv, err := NewForTesting(cfg)
	require.NoError(t, err)
	require.NotNil(t, srv)

	// Test that timeouts are correctly set
	assert.Equal(t, 30*time.Second, srv.httpSrv.ReadTimeout)
	assert.Equal(t, 30*time.Second, srv.httpSrv.WriteTimeout)
	assert.Equal(t, 120*time.Second, srv.httpSrv.IdleTimeout)
}
