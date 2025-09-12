package storage

import "fmt"

//||------------------------------------------------------------------------------------------------||
//|| Put
//||------------------------------------------------------------------------------------------------||

func (s *Storage) Put(objectName string, data []byte) error {
	if s.service == nil {
		return fmt.Errorf("storage backend not initialized")
	}
	return s.service.Put(objectName, data)
}

//||------------------------------------------------------------------------------------------------||
//|| Get
//||------------------------------------------------------------------------------------------------||

func (s *Storage) Get(objectName string) ([]byte, error) {
	if s.service == nil {
		return nil, fmt.Errorf("storage backend not initialized")
	}
	return s.service.Get(objectName)
}

//||------------------------------------------------------------------------------------------------||
//|| Delete
//||------------------------------------------------------------------------------------------------||

func (s *Storage) Delete(objectName string) error {
	if s.service == nil {
		return fmt.Errorf("storage backend not initialized")
	}
	return s.service.Delete(objectName)
}
