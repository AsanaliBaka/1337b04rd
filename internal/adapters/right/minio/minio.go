package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7"
)

type ImageStorage struct {
	client     *minio.Client
	bucketName string
}

func NewImageStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*ImageStorage, error) {
	client, err := connectMinioWithRetry(endpoint, accessKey, secretKey, useSSL, 10, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	if err := ensureBucketExists(context.Background(), client, bucketName); err != nil {
		return nil, fmt.Errorf("bucket check/create failed: %w", err)
	}

	return &ImageStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func connectMinioWithRetry(endpoint, accessKey, secretKey string, useSSL bool, retries int, delay time.Duration) (*minio.Client, error) {
	var client *minio.Client
	var err error

	for i := 0; i < retries; i++ {
		client, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})
		if err == nil {
			if _, err = client.ListBuckets(context.Background()); err == nil {
				return client, nil
			}
		}
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("could not connect to MinIO after %d retries: %w", retries, err)
}

func ensureBucketExists(ctx context.Context, client *minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	}
	return nil
}

func (u *ImageStorage) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	extension := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("post_%d%s", time.Now().UnixNano(), extension)

	_, err := u.client.PutObject(ctx, u.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return objectName, nil
}

func (u *ImageStorage) GetImage(ctx context.Context, objectName string) ([]byte, string, error) {
	object, err := u.client.GetObject(ctx, u.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object from MinIO: %w", err)
	}
	defer object.Close()

	buffer := new(bytes.Buffer)
	header := make([]byte, 512)
	n, err := object.Read(header)
	if err != nil && err != io.EOF {
		return nil, "", fmt.Errorf("failed to read object header: %w", err)
	}
	buffer.Write(header[:n])

	if _, err := io.Copy(buffer, object); err != nil {
		return nil, "", fmt.Errorf("failed to read full object: %w", err)
	}

	contentType := http.DetectContentType(header[:n])
	return buffer.Bytes(), contentType, nil
}
