package storage

import "aria/config"

//||------------------------------------------------------------------------------------------------||
//|| ConvertFromConfig: Helper to create StoreConfig from config.StorageInstanceConfig
//||------------------------------------------------------------------------------------------------||

func ConvertFromConfig(cfg config.StorageInstanceConfig) StoreConfig {
	return StoreConfig{
		Backend:   StoreBackend(cfg.Backend),
		Bucket:    cfg.Bucket,
		Region:    cfg.Region,
		Endpoint:  cfg.Endpoint,
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		LocalPath: cfg.Dir,
		// Add mapping for cloud-specific fields as needed
	}
}
