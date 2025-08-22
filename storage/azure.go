package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

//||------------------------------------------------------------------------------------------------||
//|| readSeekCloser Helper
//||------------------------------------------------------------------------------------------------||

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error { return nil }

func newReadSeekCloser(b []byte) *readSeekCloser {
	return &readSeekCloser{bytes.NewReader(b)}
}

//||------------------------------------------------------------------------------------------------||
//|| AzureBackend Struct
//||------------------------------------------------------------------------------------------------||

type AzureBackend struct {
	containerClient *container.Client
	config          StoreConfig
}

//||------------------------------------------------------------------------------------------------||
//|| NewAzureBackend Constructor
//||------------------------------------------------------------------------------------------------||

func NewAzureBackend(cfg StoreConfig) (*AzureBackend, error) {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s", cfg.AccountName, cfg.Bucket)
	cred, err := azblob.NewSharedKeyCredential(cfg.AccountName, cfg.AccountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}
	containerClient, err := container.NewClientWithSharedKeyCredential(url, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure container client: %w", err)
	}
	return &AzureBackend{
		containerClient: containerClient,
		config:          cfg,
	}, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Put: Upload an object
//||------------------------------------------------------------------------------------------------||

func (a *AzureBackend) Put(objectName string, data []byte) error {
	ctx := context.Background()
	blobClient := a.containerClient.NewBlockBlobClient(objectName)
	rsc := newReadSeekCloser(data)
	_, err := blobClient.Upload(ctx, rsc, nil)
	return err
}

//||------------------------------------------------------------------------------------------------||
//|| Get: Download an object
//||------------------------------------------------------------------------------------------------||

func (a *AzureBackend) Get(objectName string) ([]byte, error) {
	ctx := context.Background()
	blobClient := a.containerClient.NewBlockBlobClient(objectName)
	resp, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//||------------------------------------------------------------------------------------------------||
//|| Delete: Delete an object
//||------------------------------------------------------------------------------------------------||

func (a *AzureBackend) Delete(objectName string) error {
	ctx := context.Background()
	blobClient := a.containerClient.NewBlockBlobClient(objectName)
	_, err := blobClient.Delete(ctx, nil)
	return err
}
