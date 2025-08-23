package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//||------------------------------------------------------------------------------------------------||
//|| LocalBackend Struct
//||------------------------------------------------------------------------------------------------||

type StorageEngineLocal struct {
	basePath string
	config   StoreConfig
}

//||------------------------------------------------------------------------------------------------||
//|| NewLocalBackend Constructor
//||------------------------------------------------------------------------------------------------||

func NewLocalBackend(cfg StoreConfig) *StorageEngineLocal {
	path := cfg.LocalPath
	_ = os.MkdirAll(path, 0770)
	return &StorageEngineLocal{
		basePath: path,
		config:   cfg,
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Put: Write file
//||------------------------------------------------------------------------------------------------||

func (l *StorageEngineLocal) Put(objectName string, data []byte) error {
	fullPath := filepath.Join(l.basePath, objectName)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0770); err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, data, 0660)
}

//||------------------------------------------------------------------------------------------------||
//|| Get: Read file
//||------------------------------------------------------------------------------------------------||

func (l *StorageEngineLocal) Get(objectName string) ([]byte, error) {
	fullPath := filepath.Join(l.basePath, objectName)
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Delete: Remove file
//||------------------------------------------------------------------------------------------------||

func (l *StorageEngineLocal) Delete(objectName string) error {
	fullPath := filepath.Join(l.basePath, objectName)
	return os.Remove(fullPath)
}
