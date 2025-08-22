package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//||------------------------------------------------------------------------------------------------||
//|| S3Backend Struct
//||------------------------------------------------------------------------------------------------||

type S3Backend struct {
	client *s3.Client
	config StoreConfig
}

//||------------------------------------------------------------------------------------------------||
//|| NewS3Backend Constructor
//||------------------------------------------------------------------------------------------------||

func NewS3Backend(cfg StoreConfig) (*S3Backend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var awsCfg aws.Config
	var err error
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		)
	} else {
		// Uses environment/AWS IAM credentials
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg)
	return &S3Backend{
		client: client,
		config: cfg,
	}, nil
}

//||------------------------------------------------------------------------------------------------||
//|| Put: Upload an object
//||------------------------------------------------------------------------------------------------||

func (s *S3Backend) Put(objectName string, data []byte) error {
	ctx := context.Background()
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(objectName),
		Body:   bytes.NewReader(data),
	})
	return err
}

//||------------------------------------------------------------------------------------------------||
//|| Get: Download an object
//||------------------------------------------------------------------------------------------------||

func (s *S3Backend) Get(objectName string) ([]byte, error) {
	ctx := context.Background()
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

//||------------------------------------------------------------------------------------------------||
//|| Delete: Delete an object
//||------------------------------------------------------------------------------------------------||

func (s *S3Backend) Delete(objectName string) error {
	ctx := context.Background()
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(objectName),
	})
	return err
}
