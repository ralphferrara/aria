package storage

import (
	"fmt"
	"time"
)

//||------------------------------------------------------------------------------------------------||
//|| Ping (generic for any backend)
//||------------------------------------------------------------------------------------------------||

func (s *Storage) Ping() error {
	if s == nil || s.service == nil {
		return fmt.Errorf("storage service not initialized")
	}

	// Lightweight test: try to write/read/delete a dummy object
	testKey := fmt.Sprintf("healthcheck-%d", time.Now().UnixNano())
	testData := []byte("ok")

	if err := s.service.Put(testKey, testData); err != nil {
		return fmt.Errorf("put failed: %w", err)
	}
	defer s.service.Delete(testKey)

	data, err := s.service.Get(testKey)
	if err != nil {
		return fmt.Errorf("get failed: %w", err)
	}
	if string(data) != "ok" {
		return fmt.Errorf("storage data mismatch")
	}
	return nil
}
