package storage

import (
	"fmt"
	"strings"

	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| StoreService Interface
//||------------------------------------------------------------------------------------------------||

type StoreService interface {
	Put(objectName string, data []byte) error
	Get(objectName string) ([]byte, error)
	Delete(objectName string) error
}

//||------------------------------------------------------------------------------------------------||
//|| Storage Struct
//||------------------------------------------------------------------------------------------------||

type Storage struct {
	Config  StoreConfig
	service StoreService
}

//||------------------------------------------------------------------------------------------------||
//|| Global DefaultStorage
//||------------------------------------------------------------------------------------------------||

var DefaultStorage *Storage

//||------------------------------------------------------------------------------------------------||
//|| Init (Selects backend implementation)
//||------------------------------------------------------------------------------------------------||

func Init(cfg *config.Config) (map[string]*Storage, error) {
	storages := make(map[string]*Storage)
	for name, sCfg := range cfg.Storage {
		storeCfg := ConvertFromConfig(sCfg)
		st := &Storage{Config: storeCfg}

		if err := st.InitStorage(); err != nil {
			return nil, fmt.Errorf("storage '%s' init failed: %w", name, err)
		}

		fmt.Printf("\n[STRG] - Initializing storage: %s (backend: %s)", name, storeCfg.Backend)
		storages[name] = st

		if DefaultStorage == nil {
			DefaultStorage = st
		}
	}
	return storages, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Init (Selects backend implementation)
//||------------------------------------------------------------------------------------------------||

func (s *Storage) InitStorage() error {
	switch strings.ToUpper(string(s.Config.Backend)) {
	//||------------------------------------------------------------------------------------------------||
	//|| S3
	//||------------------------------------------------------------------------------------------------||
	case strings.ToUpper(string(StorageS3)):
		s3svc, err := NewS3Backend(s.Config)
		if err != nil {
			return fmt.Errorf("S3 backend init failed: %w", err)
		}
		s.service = s3svc
		return nil
	//||------------------------------------------------------------------------------------------------||
	//|| Min.io
	//||------------------------------------------------------------------------------------------------||
	case strings.ToUpper(string(StorageMinIO)):
		minioSvc, err := NewMinioBackend(s.Config)
		if err != nil {
			return fmt.Errorf("MinIO backend init failed: %w", err)
		}
		s.service = minioSvc
		return nil
	//||------------------------------------------------------------------------------------------------||
	//|| Azure
	//||------------------------------------------------------------------------------------------------||
	case strings.ToUpper(string(StorageAzure)):
		azureSvc, err := NewAzureBackend(s.Config)
		if err != nil {
			return fmt.Errorf("Azure backend init failed: %w", err)
		}
		s.service = azureSvc
		return nil
	//||------------------------------------------------------------------------------------------------||
	//|| BackendGCP
	//||------------------------------------------------------------------------------------------------||
	case strings.ToUpper(string(StorageGCP)):
		gcpSvc, err := NewGCPBackend(s.Config)
		if err != nil {
			return fmt.Errorf("GCP backend init failed: %w", err)
		}
		s.service = gcpSvc
		return nil
	//||------------------------------------------------------------------------------------------------||
	//|| Local
	//||------------------------------------------------------------------------------------------------||
	case strings.ToUpper(string(StorageLocal)):
		localSvc := NewLocalBackend(s.Config)
		s.service = localSvc
		return nil
	//||------------------------------------------------------------------------------------------------||
	//|| Fail
	//||------------------------------------------------------------------------------------------------||
	default:
		return fmt.Errorf("unsupported storage backend: %s", s.Config.Backend)
	}
}
