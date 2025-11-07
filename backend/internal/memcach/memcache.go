package memcach

import (
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
}

func (m *MemCache) cleanup() {

	for {
		time.Sleep(30 * time.Second)
		m.mu.Lock()

		for k, v := range m.data {
			if time.Now().UnixNano() > v.expiration {
				delete(m.data, k)
			}
		}
		m.mu.Unlock()

	}

}

// NewMemCache creates a new in-memory cache
func NewMemCache() *MemCache {

	mc := &MemCache{
		data: make(map[string]item),
	}
	go mc.cleanup()
	return mc
}

func (mc *MemCache) Set(key string, value []byte) {
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
