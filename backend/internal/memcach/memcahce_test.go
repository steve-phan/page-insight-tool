package memcach

import (
	"sync"
	"testing"
	"time"

	"github.com/steve-phan/page-insight-tool/internal/config"
)

// Test constants for reuse (best practice: avoid magic numbers)
const (
	testKey   = "testKey"
	testValue = "testValue"
)

var mockConfig = config.CacheConfig{
	Enabled:         true,
	DefaultTTL:      300 * time.Second, // Use time.Duration directly for clarity
	MaxSize:         100,
	CleanupInterval: 60 * time.Second,
}

func TestMemCache_SetAndGet(t *testing.T) {
	mc := NewMemCache(mockConfig)
	defer mc.Stop()

	mc.Set(testKey, []byte(testValue))

	retrieved, found := mc.Get(testKey)
	if !found {
		t.Fatalf("Expected to find key %s", testKey)
	}
	if string(retrieved) != testValue {
		t.Fatalf("Expected value %s, got %s", testValue, string(retrieved))
	}
}

func TestMemCache_Expiration(t *testing.T) {
	mc := NewMemCache(mockConfig)
	defer mc.Stop()

	mc.Set(testKey, []byte(testValue))
	time.Sleep(1 * time.Second) // Short sleep for quick test

	_, found := mc.Get(testKey)
	if !found {
		t.Error("Expected key to still exist before expiration")
	}

	// For full expiration test, set short TTL in config or override
}

func TestMemCache_RaceCondition(t *testing.T) {
	mc := NewMemCache(mockConfig)
	defer mc.Stop()

	var wg sync.WaitGroup

	// Writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ { // Fewer iterations for faster tests
			mc.Set("raceKey", []byte("raceValue"))
		}
	}()

	// Reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			mc.Get("raceKey")
		}
	}()

	wg.Wait()
}

// Table-driven test for multiple configs (best practice for variations)
func TestMemCache_WithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name   string
		config config.CacheConfig
	}{
		{"Default Config", mockConfig},
		{"Short TTL", config.CacheConfig{Enabled: true, DefaultTTL: 1 * time.Second, MaxSize: 10, CleanupInterval: 5 * time.Second}},
		{"Large Size", config.CacheConfig{Enabled: true, DefaultTTL: 300 * time.Second, MaxSize: 10000, CleanupInterval: 30 * time.Second}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := NewMemCache(tt.config)
			defer mc.Stop()

			mc.Set(testKey, []byte(testValue))
			_, found := mc.Get(testKey)
			if !found {
				t.Errorf("Failed to set/get with config: %+v", tt.config)
			}
		})
	}
}
