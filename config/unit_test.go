package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndGetConfig(t *testing.T) {
	// Locate test config file (adjust path as needed)
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, "..", "testdata", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("Test config file not found: %s", configPath)
	}

	cfg, err := Init(configPath)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Init returned nil config")
	}

	// Verify GetConfig returns the same pointer
	got := GetConfig()
	if got != cfg {
		t.Error("GetConfig did not return the loaded config pointer")
	}

	// Spot check some fields
	if cfg.App.Name == "" {
		t.Error("App.Name is empty")
	}
	if len(cfg.DB) == 0 {
		t.Error("No DB config found")
	}
	if len(cfg.Cache) == 0 {
		t.Error("No Cache config found")
	}
	if len(cfg.Storage) == 0 {
		t.Error("No Storage config found")
	}
	if len(cfg.Queue) == 0 {
		t.Error("No Queue config found")
	}
	if cfg.Auth.JwtSecret == "" {
		t.Error("Auth.JwtSecret is empty")
	}
	if cfg.Locale.Default == "" {
		t.Error("Locale.Default is empty")
	}
	if cfg.Template.Dir == "" {
		t.Error("Template.Dir is empty")
	}
}
