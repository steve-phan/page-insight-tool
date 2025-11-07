package memcach

import (
	"fmt"
	"sync"
	"time"
)

// Implement a basic MemCache

const memcacheDefaultExpiration = 5 * 60 // 5 minutes in seconds

type item struct {
	value      []byte
	expiration int64
}
type MemCache struct {
	data map[string]item
	mu   sync.RWMutex

	stop chan struct{}
}

func (m *MemCache) cleanup() {

	// Create a ticker that ticks every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
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
			fmt.Println("Greatful stop interval clean up")
			return
		}
	}
}

func (m *MemCache) Stop() {
	close(m.stop)
}

// NewMemCache creates a new in-memory cache
func NewMemCache() *MemCache {

	mc := &MemCache{
		data: make(map[string]item),
		stop: make(chan struct{}),
	}
	go mc.cleanup()
	return mc
}

func (mc *MemCache) Set(key string, value []byte) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data[key] = item{
		value:      value,
		expiration: int64(time.Now().Add(memcacheDefaultExpiration * time.Second).UnixNano()),
	}
}

// Get retrieves an item from the cache
func (mc *MemCache) Get(key string) ([]byte, bool) {
	itm, found := mc.data[key]
	if !found || time.Now().UnixNano() > itm.expiration {
		return nil, false
	}
	return itm.value, true
}

// Get Memcache interface
type Memcache interface {
	Set(key string, value []byte)
	Get(key string) ([]byte, bool)
	Stop()
}

var globalCache *MemCache
var once sync.Once

// GetMemcache returns the global Memcache instance
func GetMemcache() Memcache {
	once.Do(func() {
		globalCache = NewMemCache()
	})

	return globalCache
}
