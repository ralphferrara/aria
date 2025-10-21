package storage

import "github.com/ralphferrara/aria/config"

//||------------------------------------------------------------------------------------------------||
//|| ConvertFromConfig: Helper to create StoreConfig from config.StorageInstanceConfig
//||------------------------------------------------------------------------------------------------||

func ConvertFromConfig(cfg config.StorageInstanceConfig) StoreConfig {
	return StoreConfig{
		Backend:         StoreBackend(cfg.Backend),
		Bucket:          cfg.Bucket,
		Region:          cfg.Region,
		CredentialsJSON: cfg.CredentialsJSON,
		Project:         cfg.Project,
		Endpoint:        cfg.Endpoint,
		AccessKey:       cfg.AccessKey,
		SecretKey:       cfg.SecretKey,
		LocalPath:       cfg.Dir,
	}
}
