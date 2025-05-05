package right

import (
	"time"

	"1337b04rd/internal/domain"
)

type PostRepository interface {
	CreatePost(post *domain.Post) error
	GetPostByID(id string) (*domain.Post, error)
	GetPosts() ([]*domain.Post, error)
	CreateComment(comment *domain.Comment) error
	DeleteOldPosts(noCommentsThreshold, withCommentsThreshold time.Duration) error
	Close() error
}
