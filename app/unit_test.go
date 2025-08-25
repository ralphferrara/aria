package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	// Locate test config (adjust if needed)
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, "..", "testdata", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("Test config file not found: %s", configPath)
	}

	err := Init(configPath)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Basic assertions
	if Config == nil {
		t.Error("Config not set after Init")
	}
	if len(Storages) == 0 {
		t.Error("No storages initialized")
	}
	if len(SQLDB) == 0 && len(MongoDB) == 0 {
		t.Error("No databases initialized")
	}
	if len(QueueRabbit) == 0 {
		t.Error("No queues initialized")
	}
	if len(CacheRedis) == 0 && len(CacheKeyDB) == 0 && len(CacheMemcached) == 0 && len(CacheMemory) == 0 {
		t.Error("No caches initialized")
	}
}
