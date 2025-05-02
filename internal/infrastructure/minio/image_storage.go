package minio

import (
	. "1337b04rd/internal/domain"
	"1337b04rd/internal/infrastructure/config"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type imageStorage struct {
	client     *minio.Client
	bucketName string
}

const ImageURLExpiration = 7 * 24 * time.Hour

func NewImageStrorage(cfg *config.Config, ctx context.Context) (ImageStorage, error) {
	client, err := minio.New(cfg.MinioEndPoint, &minio.Options{

		Creds:  credentials.NewStaticV4(cfg.MinioUser, cfg.MinioPassword, ""),
		Secure: cfg.MinioSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.BucketName)

	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &imageStorage{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

func (i *imageStorage) CreateImage(ctx context.Context, imageName string, imageData io.Reader, size int64) (string, error) {
	_, err := i.client.PutObject(ctx, i.bucketName, imageName, imageData, size, minio.PutObjectOptions{})

	if err != nil {
		return "", fmt.Errorf("failed to create image: %w", err)
	}

	return imageName, nil
}

func (i *imageStorage) GetImageURL(ctx context.Context, imageName string) (string, error) {
	url, err := i.client.PresignedGetObject(ctx, i.bucketName, imageName, ImageURLExpiration, nil)

	if err != nil {
		return "", fmt.Errorf("failed to generate image URL: %w", err)
	}

	return url.String(), nil
}
func (i *imageStorage) DeleteImage(ctx context.Context, imageName string) error {
	err := i.client.RemoveObject(ctx, i.bucketName, imageName, minio.RemoveObjectOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}
