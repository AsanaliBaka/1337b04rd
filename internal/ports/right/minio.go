package right

import "io"

type ImageStorage interface {
	UploadImage(bucketName string, fileName string, imageData io.Reader, size int64) (string, error)
	GetImageURL(bucketName string, fileName string) (string, error)
	CreateBucket(bucketName string) error
	BucketExists(bucketName string) (bool, error)
}
