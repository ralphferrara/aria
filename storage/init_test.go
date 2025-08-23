//||------------------------------------------------------------------------------------------------||
//|| Storage Backend Tests
//||------------------------------------------------------------------------------------------------||

package storage

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

//||------------------------------------------------------------------------------------------------||
//|| Helper: Test Data
//||------------------------------------------------------------------------------------------------||

var (
	testKey   = "testfile.txt"
	testValue = []byte("Go storage integration test value!")
)

//||------------------------------------------------------------------------------------------------||
//|| Table-driven Test
//||------------------------------------------------------------------------------------------------||

func TestStorageBackends(t *testing.T) {

	//||------------------------------------------------------------------------------------------------||
	//|| Load Env
	//||------------------------------------------------------------------------------------------------||

	_ = godotenv.Load("../.env")

	//||------------------------------------------------------------------------------------------------||
	//|| Define Test Cases
	//||------------------------------------------------------------------------------------------------||
	cases := []struct {
		name   string
		config StoreConfig
		skip   bool
	}{
		{
			name: "S3",
			config: StoreConfig{
				Backend:   StorageS3,
				Bucket:    GetEnv("STORAGE_BUCKET", ""),
				Region:    GetEnv("STORAGE_REGION", ""),
				AccessKey: GetEnv("STORAGE_ACCESS_KEY", ""),
				SecretKey: GetEnv("STORAGE_SECRET_KEY", ""),
			},
			skip: false,
		},
		{
			name: "MinIO",
			config: StoreConfig{
				Backend:   StorageMinIO,
				Bucket:    GetEnv("STORAGE_BUCKET", ""),
				Region:    GetEnv("STORAGE_REGION", ""),
				Endpoint:  GetEnv("STORAGE_ENDPOINT", ""),
				AccessKey: GetEnv("STORAGE_ACCESS_KEY", ""),
				SecretKey: GetEnv("STORAGE_SECRET_KEY", ""),
				UseSSL:    GetEnvBool("STORAGE_USE_SSL", false),
			},
			skip: false,
		},
		{
			name: "Azure",
			config: StoreConfig{
				Backend:     StorageAzure,
				Bucket:      GetEnv("STORAGE_BUCKET", ""),
				AccountName: GetEnv("STORAGE_ACCOUNT_NAME", ""),
				AccountKey:  GetEnv("STORAGE_ACCOUNT_KEY", ""),
			},
			skip: false,
		},
		{
			name: "GCP",
			config: StoreConfig{
				Backend:         StorageGCP,
				Bucket:          GetEnv("STORAGE_BUCKET", ""),
				CredentialsJSON: GetEnv("STORAGE_CREDENTIALS_JSON", ""),
			},
			skip: false,
		},
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

	//||------------------------------------------------------------------------------------------------||
	//|| Run Each Backend Test
	//||------------------------------------------------------------------------------------------------||
	for _, tc := range cases {
		tc := tc // capture range variable

		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("skipped %s backend", tc.name)
			}

			// Optional: skip test if required config is missing
			if tc.name != "Local" && missingConfig(tc.config) {
				t.Skipf("missing credentials/config for %s backend", tc.name)
			}

			//||------------------------------------------------------------------------------------------------||
			//|| Init Storage
			//||------------------------------------------------------------------------------------------------||
			storage := &Storage{Config: tc.config}
			if err := storage.Init(); err != nil {
				t.Fatalf("[%s] Init failed: %v", tc.name, err)
			}

			//||------------------------------------------------------------------------------------------------||
			//|| Put
			//||------------------------------------------------------------------------------------------------||
			if err := storage.Put(testKey, testValue); err != nil {
				t.Fatalf("[%s] Put failed: %v", tc.name, err)
			}
			t.Logf("[%s] Put success!", tc.name)

			//||------------------------------------------------------------------------------------------------||
			//|| Get
			//||------------------------------------------------------------------------------------------------||
			got, err := storage.Get(testKey)
			if err != nil {
				t.Fatalf("[%s] Get failed: %v", tc.name, err)
			}
			if string(got) != string(testValue) {
				t.Fatalf("[%s] Get returned wrong value: got %q want %q", tc.name, got, testValue)
			}
			t.Logf("[%s] Get success!", tc.name)

			//||------------------------------------------------------------------------------------------------||
			//|| Delete
			//||------------------------------------------------------------------------------------------------||
			if err := storage.Delete(testKey); err != nil {
				t.Fatalf("[%s] Delete failed: %v", tc.name, err)
			}
			t.Logf("[%s] Delete success!", tc.name)

			//||------------------------------------------------------------------------------------------------||
			//|| Cleanup (local)
			//||------------------------------------------------------------------------------------------------||
			if tc.name == "Local" {
				os.RemoveAll(tc.config.LocalPath)
			}
		})
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Helper: Check if Any Essential Config Is Missing
//||------------------------------------------------------------------------------------------------||

func missingConfig(cfg StoreConfig) bool {
	switch cfg.Backend {
	case StorageS3:
		return cfg.Bucket == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Region == ""
	case StorageMinIO:
		return cfg.Bucket == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == ""
	case StorageAzure:
		return cfg.Bucket == "" || cfg.AccountName == "" || cfg.AccountKey == ""
	case StorageGCP:
		return cfg.Bucket == "" || cfg.CredentialsJSON == ""
	}
	return false
}

//||------------------------------------------------------------------------------------------------||
//|| Get Env with Default
//||------------------------------------------------------------------------------------------------||

func GetEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

//||------------------------------------------------------------------------------------------------||
//|| Check if Boolean Env is True
//||------------------------------------------------------------------------------------------------||

func GetEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "true" || val == "1" {
		return true
	}
	if val == "false" || val == "0" {
		return false
	}
	return def
}
