//||------------------------------------------------------------------------------------------------||
//|| Storage Backend Tests
//||------------------------------------------------------------------------------------------------||

package storage

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"base/helpers"
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
				Bucket:    helpers.GetEnv("STORAGE_BUCKET", ""),
				Region:    helpers.GetEnv("STORAGE_REGION", ""),
				AccessKey: helpers.GetEnv("STORAGE_ACCESS_KEY", ""),
				SecretKey: helpers.GetEnv("STORAGE_SECRET_KEY", ""),
			},
			skip: false,
		},
		{
			name: "MinIO",
			config: StoreConfig{
				Backend:   StorageMinIO,
				Bucket:    helpers.GetEnv("STORAGE_BUCKET", ""),
				Region:    helpers.GetEnv("STORAGE_REGION", ""),
				Endpoint:  helpers.GetEnv("STORAGE_ENDPOINT", ""),
				AccessKey: helpers.GetEnv("STORAGE_ACCESS_KEY", ""),
				SecretKey: helpers.GetEnv("STORAGE_SECRET_KEY", ""),
				UseSSL:    helpers.GetEnvBool("STORAGE_USE_SSL", false),
			},
			skip: false,
		},
		{
			name: "Azure",
			config: StoreConfig{
				Backend:     StorageAzure,
				Bucket:      helpers.GetEnv("STORAGE_BUCKET", ""),
				AccountName: helpers.GetEnv("STORAGE_ACCOUNT_NAME", ""),
				AccountKey:  helpers.GetEnv("STORAGE_ACCOUNT_KEY", ""),
			},
			skip: false,
		},
		{
			name: "GCP",
			config: StoreConfig{
				Backend:         StorageGCP,
				Bucket:          helpers.GetEnv("STORAGE_BUCKET", ""),
				CredentialsJSON: helpers.GetEnv("STORAGE_CREDENTIALS_JSON", ""),
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
