package memcach

import (
	"sync"
	"testing"
)

func TestMemCache_SetAndGet(t *testing.T) {

	mc := NewMemCache()
	key := "testKey"
	value := []byte("testValue")
	mc.Set(key, value)

	retrieved, found := mc.Get(key)
	if !found {
		t.Fatalf("Expected to find key %s", key)
	}
	if string(retrieved) != string(value) {
		t.Fatalf("Expected value %s, got %s", value, retrieved)
	}

}

func TestMemCache_RaceCondition(t *testing.T) {
	mc := NewMemCache()
	defer mc.Stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			mc.Set("raceKey", []byte("raceValue"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			mc.Get("raceKey")
		}
	}()

	wg.Wait()
}
