package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/ralphferrara/aria/log"
	"google.golang.org/api/option"
)

//||------------------------------------------------------------------------------------------------||
//|| GCPBackend Struct
//||------------------------------------------------------------------------------------------------||

type StorageEngineGCP struct {
	client *storage.Client
	bucket string
	config StoreConfig
}

//||------------------------------------------------------------------------------------------------||
//|| NewGCPBackend Constructor
//||------------------------------------------------------------------------------------------------||

func NewGCPBackend(cfg StoreConfig) (*StorageEngineGCP, error) {
	ctx := context.Background()
	var client *storage.Client
	var err error

	log.PrettyPrint(cfg)

	if cfg.CredentialsJSON != "" {
		fmt.Println("[GCP] - Credentials Path:", cfg.CredentialsJSON)
		if _, err := os.Stat(cfg.CredentialsJSON); err != nil {
			fmt.Println("[GCP] - File not found:", err)
		}
		if _, statErr := os.Stat(cfg.CredentialsJSON); statErr == nil {
			absPath, _ := filepath.Abs(cfg.CredentialsJSON)
			creds, readErr := os.ReadFile(absPath)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read GCP credentials file: %w", readErr)
			}
			client, err = storage.NewClient(ctx, option.WithCredentialsJSON(creds))
		} else {
			// Otherwise, assume it's already the JSON content
			client, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(cfg.CredentialsJSON)))
		}
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCP Storage client: %w", err)
	}

	return &StorageEngineGCP{
		client: client,
		bucket: cfg.Bucket,
		config: cfg,
	}, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Put: Upload an object
//||------------------------------------------------------------------------------------------------||

func (g *StorageEngineGCP) Put(objectName string, data []byte) error {
	ctx := context.Background()
	wc := g.client.Bucket(g.bucket).Object(objectName).NewWriter(ctx)
	_, err := wc.Write(data)
	if err != nil {
		_ = wc.Close()
		return err
	}
	return wc.Close()
}

//||------------------------------------------------------------------------------------------------||
//|| Get: Download an object
//||------------------------------------------------------------------------------------------------||

func (g *StorageEngineGCP) Get(objectName string) ([]byte, error) {
	ctx := context.Background()
	rc, err := g.client.Bucket(g.bucket).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, rc)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//||------------------------------------------------------------------------------------------------||
//|| Delete: Delete an object
//||------------------------------------------------------------------------------------------------||

func (g *StorageEngineGCP) Delete(objectName string) error {
	ctx := context.Background()
	return g.client.Bucket(g.bucket).Object(objectName).Delete(ctx)
}
