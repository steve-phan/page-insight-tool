package config

import (
	"os"
	"testing"
	"time"
)

const configPath = "../../config/config.yaml"

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{
			name:        "valid config file",
			configPath:  configPath,
			expectError: false,
		},
		{
			name:        "non-existent config file",
			configPath:  "non-existent.yaml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(tt.configPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if cfg == nil {
				t.Errorf("config should not be nil")
				return
			}

			// Test default values
			if cfg.Server.Port != 8080 {
				t.Errorf("expected port 8080, got %d", cfg.Server.Port)
			}

			if cfg.App.Name != "Page Insight Tool" {
				t.Errorf("expected app name 'Page Insight Tool', got '%s'", cfg.App.Name)
			}
		})
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("PIT_SERVER_PORT", "9090")
	os.Setenv("PIT_APP_NAME", "Test App")
	defer func() {
		os.Unsetenv("PIT_SERVER_PORT")
		os.Unsetenv("PIT_APP_NAME")
	}()

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Test environment overrides
	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090 from env, got %d", cfg.Server.Port)
	}

	if cfg.App.Name != "Test App" {
		t.Errorf("expected app name 'Test App' from env, got '%s'", cfg.App.Name)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{
					Port: 8080,
				},
			},
			expectError: false,
		},
		{
			name: "invalid port - too low",
			config: &Config{
				Server: ServerConfig{
					Port: 0,
				},
			},
			expectError: true,
		},
		{
			name: "invalid port - too high",
			config: &Config{
				Server: ServerConfig{
					Port: 70000,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestConfigMethods(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		App: AppConfig{
			Environment: "production",
		},
	}

	// Test GetAddress
	expected := "localhost:8080"
	if got := cfg.GetAddress(); got != expected {
		t.Errorf("GetAddress() = %v, want %v", got, expected)
	}

	// Test IsProduction
	if !cfg.IsProduction() {
		t.Errorf("IsProduction() = false, want true")
	}

	// Test IsDevelopment
	cfg.App.Environment = "development"
	if !cfg.IsDevelopment() {
		t.Errorf("IsDevelopment() = false, want true")
	}
}

func TestDurationParsing(t *testing.T) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Test that durations are parsed correctly
	expectedReadTimeout := 30 * time.Second
	if cfg.Server.ReadTimeout != expectedReadTimeout {
		t.Errorf("ReadTimeout = %v, want %v", cfg.Server.ReadTimeout, expectedReadTimeout)
	}

	expectedWriteTimeout := 30 * time.Second
	if cfg.Server.WriteTimeout != expectedWriteTimeout {
		t.Errorf("WriteTimeout = %v, want %v", cfg.Server.WriteTimeout, expectedWriteTimeout)
	}
}
