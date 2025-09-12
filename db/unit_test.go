package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ralphferrara/aria/config"
)

func TestDBInit(t *testing.T) {
	// Locate test config (adjust if needed)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}
	configPath := filepath.Join(cwd, "..", "testdata", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("Test config file not found: %s", configPath)
	}

	cfg, err := config.Init(configPath)
	if err != nil {
		t.Fatalf("config.Init failed: %v", err)
	}

	sqlMap, mongoMap, err := Init(cfg)
	if err != nil {
		t.Fatalf("db.Init failed: %v", err)
	}

	SQL = sqlMap
	Mongo = mongoMap

	if len(SQL) == 0 && len(Mongo) == 0 {
		t.Error("No DB connections were initialized (SQL and Mongo maps are empty)")
	}

	// Optionally, check each SQL connection
	for name, wrapper := range SQL {
		sqlDB, err := wrapper.DB.DB()
		if err != nil {
			t.Errorf("GORM DB.DB() failed for %s: %v", name, err)
			continue
		}
		if err := sqlDB.Ping(); err != nil {
			t.Errorf("SQL DB ping failed for %s: %v", name, err)
		}
	}

	// Optionally, check each Mongo connection (ping)
	for name, wrapper := range Mongo {
		if err := wrapper.Database.Client().Ping(nil, nil); err != nil {
			t.Errorf("Mongo DB ping failed for %s: %v", name, err)
		}
	}
}
