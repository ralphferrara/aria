package storage

import (
	"fmt"
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

func (s *Storage) Init() error {
	switch s.Config.Backend {
	case BackendS3:
		s3svc, err := NewS3Backend(s.Config)
		if err != nil {
			return fmt.Errorf("S3 backend init failed: %w", err)
		}
		s.service = s3svc
		return nil
	case BackendMinIO:
		minioSvc, err := NewMinioBackend(s.Config)
		if err != nil {
			return fmt.Errorf("MinIO backend init failed: %w", err)
		}
		s.service = minioSvc
		return nil
	case BackendAzure:
		azureSvc, err := NewAzureBackend(s.Config)
		if err != nil {
			return fmt.Errorf("Azure backend init failed: %w", err)
		}
		s.service = azureSvc
		return nil
	case BackendGCP:
		gcpSvc, err := NewGCPBackend(s.Config)
		if err != nil {
			return fmt.Errorf("GCP backend init failed: %w", err)
		}
		s.service = gcpSvc
		return nil
	case BackendLocal:
		localSvc := NewLocalBackend(s.Config)
		s.service = localSvc
		return nil
	default:
		return fmt.Errorf("unsupported storage backend: %s", s.Config.Backend)
	}
}
