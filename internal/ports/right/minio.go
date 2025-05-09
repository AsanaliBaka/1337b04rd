package right

import (
	"context"
	"mime/multipart"
)

type MinioPort interface {
	ImageStorage
}
type ImageStorage interface {
	UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	GetImage(ctx context.Context, imageName string) ([]byte, string, error)
}
