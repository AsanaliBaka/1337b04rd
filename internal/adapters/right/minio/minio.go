package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint string
	User     string
	Password string
	UseSSL   bool
}

type MinIOClient struct {
	client   *minio.Client
	endpoint string
}

func DefaultConfig() Config {
	return Config{
		Endpoint: "localhost:9000",
		User:     "minioadmin",
		Password: "minioadmin",
		UseSSL:   false,
	}
}

func New(cfg Config) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	m := &MinIOClient{
		client:   client,
		endpoint: cfg.Endpoint,
	}

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.ListBuckets(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	// Создаем необходимые бакеты
	for _, bucket := range []string{"images", "thumbnails"} {
		if err := m.CreateBucket(ctx, bucket); err != nil {
			return nil, fmt.Errorf("failed to initialize bucket %s: %w", bucket, err)
		}

		// Настраиваем политику доступа
		policy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`

		if err := client.SetBucketPolicy(ctx, bucket, fmt.Sprintf(policy, bucket)); err != nil {
			slog.Warn("Failed to set bucket policy", "bucket", bucket, "error", err)
		}
	}

	return m, nil
}

func (m *MinIOClient) CreateBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("bucket check failed: %w", err)
	}
	if !exists {
		if err := m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("bucket creation failed: %w", err)
		}
	}
	return nil
}

func (m *MinIOClient) UploadImage(ctx context.Context, bucket, objectName string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := m.client.PutObject(ctx, bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	slog.Info("Image uploaded successfully",
		"bucket", bucket,
		"object", objectName,
		"size", len(data),
	)

	return m.GetObjectURL(ctx, bucket, objectName)
}

func (m *MinIOClient) DownloadObject(ctx context.Context, bucket, objectName string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer object.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, object); err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return buf.Bytes(), nil
}

func (m *MinIOClient) DeleteObject(ctx context.Context, bucket, objectName string) error {
	if err := m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

func (m *MinIOClient) ListObjects(ctx context.Context, bucket string) ([]string, error) {
	var objectURLs []string

	objectCh := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			slog.Warn("Error listing object", "error", object.Err)
			continue
		}

		url, err := m.GetObjectURL(ctx, bucket, object.Key)
		if err != nil {
			slog.Warn("Error generating URL", "error", err)
			continue
		}

		objectURLs = append(objectURLs, url)
	}

	return objectURLs, nil
}

func (m *MinIOClient) GetObjectURL(ctx context.Context, bucket, objectName string) (string, error) {
	reqParams := make(url.Values)
	// Генерируем URL с временным токеном (действителен 24 часа)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucket, objectName, 24*time.Hour, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}
