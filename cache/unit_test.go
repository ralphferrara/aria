//||------------------------------------------------------------------------------------------------||
//|| Cache Package: Tests
//|| cache_test.go
//||------------------------------------------------------------------------------------------------||

package cache

import (
	"testing"
	"time"

	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Helper: Init Test Config
//||------------------------------------------------------------------------------------------------||

func initTestConfig(t *testing.T) *config.Config {
	cfg, err := config.Init("../config.sample.json")
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
	return cfg
}

//||------------------------------------------------------------------------------------------------||
//|| Test: Cache Backends
//||------------------------------------------------------------------------------------------------||

func TestCacheBackends(t *testing.T) {
	// Load config and init all cache backends
	_ = initTestConfig(t)
	if err := Init(); err != nil {
		t.Fatalf("cache init failed: %v", err)
	}

	// Test Redis/KeyDB (if defined)
	for name, c := range Redis {
		t.Run("Redis_"+name, func(t *testing.T) {
			key, val := "unit_test_key", "redis_works"
			err := c.Set(key, val, 2*time.Second)
			if err != nil {
				t.Fatalf("redis set failed: %v", err)
			}
			got, err := c.Get(key)
			if err != nil {
				t.Fatalf("redis get failed: %v", err)
			}
			if got != val {
				t.Fatalf("redis get mismatch: want %q, got %q", val, got)
			}
		})
	}
	for name, c := range KeyDB {
		t.Run("KeyDB_"+name, func(t *testing.T) {
			key, val := "unit_test_key", "keydb_works"
			err := c.Set(key, val, 2*time.Second)
			if err != nil {
				t.Fatalf("keydb set failed: %v", err)
			}
			got, err := c.Get(key)
			if err != nil {
				t.Fatalf("keydb get failed: %v", err)
			}
			if got != val {
				t.Fatalf("keydb get mismatch: want %q, got %q", val, got)
			}
		})
	}
	for name, c := range Memcached {
		t.Run("Memcached_"+name, func(t *testing.T) {
			key, val := "unit_test_key", []byte("memcached_works")
			err := c.Set(key, val, 2)
			if err != nil {
				t.Fatalf("memcached set failed: %v", err)
			}
			got, err := c.Get(key)
			if err != nil {
				t.Fatalf("memcached get failed: %v", err)
			}
			if string(got) != string(val) {
				t.Fatalf("memcached get mismatch: want %q, got %q", val, got)
			}
		})
	}
	for name, c := range Memory {
		t.Run("Memory_"+name, func(t *testing.T) {
			key, val := "unit_test_key", "memory_works"
			c.Set(key, val)
			got, ok := c.Get(key)
			if !ok {
				t.Fatalf("memory get not found")
			}
			if got != val {
				t.Fatalf("memory get mismatch: want %q, got %q", val, got)
			}
		})
	}
}
