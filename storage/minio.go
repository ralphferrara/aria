package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//||------------------------------------------------------------------------------------------------||
//|| MinioBackend Struct
//||------------------------------------------------------------------------------------------------||

type StorageEngineMinio struct {
	client *minio.Client
	config StoreConfig
}

//||------------------------------------------------------------------------------------------------||
//|| NewMinioBackend Constructor
//||------------------------------------------------------------------------------------------------||

func NewMinioBackend(cfg StoreConfig) (*StorageEngineMinio, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}
	return &StorageEngineMinio{
		client: client,
		config: cfg,
	}, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Put: Upload an object
//||------------------------------------------------------------------------------------------------||

func (m *StorageEngineMinio) Put(objectName string, data []byte) error {
	ctx := context.Background()
	_, err := m.client.PutObject(
		ctx,
		m.config.Bucket,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	return err
}

//||------------------------------------------------------------------------------------------------||
//|| Get: Download an object
//||------------------------------------------------------------------------------------------------||

func (m *StorageEngineMinio) Get(objectName string) ([]byte, error) {
	ctx := context.Background()
	obj, err := m.client.GetObject(ctx, m.config.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//||------------------------------------------------------------------------------------------------||
//|| Delete: Delete an object
//||------------------------------------------------------------------------------------------------||

func (m *StorageEngineMinio) Delete(objectName string) error {
	ctx := context.Background()
	return m.client.RemoveObject(ctx, m.config.Bucket, objectName, minio.RemoveObjectOptions{})
}
