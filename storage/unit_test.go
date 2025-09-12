//||------------------------------------------------------------------------------------------------||
//|| Storage Backend Tests
//|| unit_test.go
//||------------------------------------------------------------------------------------------------||

package storage

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Helper: Test Data
//||------------------------------------------------------------------------------------------------||

var (
	testKey   = "testfile.txt"
	testValue = []byte("Go storage integration test value!")
)

//||------------------------------------------------------------------------------------------------||
//|| Table-driven Test (single backends)
//||------------------------------------------------------------------------------------------------||

func TestStorageBackends(t *testing.T) {

	_ = godotenv.Load("../.env")

	cases := []struct {
		name   string
		config StoreConfig
		skip   bool
	}{
		{
			name: "Local",
			config: StoreConfig{
				Backend:   StorageLocal,
				Bucket:    "dummy",
				LocalPath: "./test_localdata",
			},
			skip: false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("skipped %s backend", tc.name)
			}

			st := &Storage{Config: tc.config}
			if err := st.InitStorage(); err != nil {
				t.Fatalf("[%s] InitStorage failed: %v", tc.name, err)
			}

			if err := st.Put(testKey, testValue); err != nil {
				t.Fatalf("[%s] Put failed: %v", tc.name, err)
			}

			got, err := st.Get(testKey)
			if err != nil {
				t.Fatalf("[%s] Get failed: %v", tc.name, err)
			}
			if string(got) != string(testValue) {
				t.Fatalf("[%s] Get returned wrong value: got %q want %q", tc.name, got, testValue)
			}

			if err := st.Delete(testKey); err != nil {
				t.Fatalf("[%s] Delete failed: %v", tc.name, err)
			}

			if tc.name == "Local" {
				os.RemoveAll(tc.config.LocalPath)
			}
		})
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Integration Test: Use main config.Init
//||------------------------------------------------------------------------------------------------||

func TestInitFromConfig(t *testing.T) {
	_ = godotenv.Load("../.env")

	configPath := filepath.Join("..", "testdata", "config.json")
	cfg, err := config.Init(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	stores, err := Init(cfg)
	if err != nil {
		t.Fatalf("Init(cfg) failed: %v", err)
	}

	if len(stores) == 0 {
		t.Fatal("expected at least one storage backend from config")
	}

	for name, st := range stores {
		t.Logf("Initialized storage: %s -> %+v", name, st.Config.Backend)
	}
}
