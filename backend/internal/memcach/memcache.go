package memcach

import (
	"fmt"
	"sync"
	"time"

	"github.com/steve-phan/page-insight-tool/internal/config"
)

const memcacheDefaultExpiration = 5 * 60 // 5 minutes in seconds

// item represents a cached item with value and expireration time
type item struct {
	value      []byte
	expiration int64
}

// InMemoryCache defines the interface for an in-memory cache operations
type InMemoryCache interface {
	Set(key string, value []byte)
	Get(key string) ([]byte, bool)
	Stop()
}

// MemCache is an in-memory cache implementation with expiration
type MemCache struct {
	config config.CacheConfig
	data   map[string]item
	mu     sync.RWMutex

	stop chan struct{}
}

var (
	globalCache InMemoryCache
	once        sync.Once
)

// NewMemCache creates a new in-memory cache instance
func NewMemCache(cacheConfig config.CacheConfig) *MemCache {

	mc := &MemCache{
		data:   make(map[string]item),
		stop:   make(chan struct{}),
		config: cacheConfig,
	}
	go mc.cleanup(cacheConfig.CleanupInterval)
	return mc
}

func InitCache(cacheConfig config.CacheConfig) {
	once.Do(func() {
		globalCache = NewMemCache(cacheConfig)
	})
}

// GetMemCache returns the global Memcache instance
func GetMemCache() InMemoryCache {

	return globalCache
}

func (mc *MemCache) Set(key string, value []byte) {
	if !mc.config.Enabled {
		return
	}
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data[key] = item{
		value:      value,
		expiration: int64(time.Now().Add(memcacheDefaultExpiration * time.Second).UnixNano()),
	}
}

// Get retrieves an item from the cache
func (mc *MemCache) Get(key string) ([]byte, bool) {
	if !mc.config.Enabled {
		return nil, false
	}

	mc.mu.RLock()
	defer mc.mu.RUnlock()
	itm, found := mc.data[key]
	if !found || time.Now().UnixNano() > itm.expiration {
		return nil, false
	}
	return itm.value, true
}

// Stop gracefully stop the MemCache cleanup goroutine
func (m *MemCache) Stop() {
	close(m.stop)
}

// cleanup periodically removes expired items from the cache
func (m *MemCache) cleanup(interval time.Duration) {

	// Create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			for k, v := range m.data {
				if time.Now().UnixNano() > v.expiration {
					delete(m.data, k)
				}

			}
			m.mu.Unlock()
		case <-m.stop: // Listen to the stop signal channel
			fmt.Println("Gracefully stopping interval cleanup")
			return
		}
	}
}
