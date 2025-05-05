package left

import "1337b04rd/internal/domain"

type PostRepository interface {
	SavePost(post *domain.Post) error
	GetPost(id string) (*domain.Post, error)
	GetPosts() ([]*domain.Post, error)
	SaveComment(comment *domain.Comment) error
	DeletePost(id string) error
}

type AvatarProvider interface {
	GetRandomAvatar() (*domain.User, error)
}

type ImageStorage interface {
	UploadImage(data []byte) (string, error)
	GetImageURL(key string) string
}
