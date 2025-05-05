package left

import "1337b04rd/internal/domain"

type PostService interface {
	CreatePost(title, content string, image []byte) (*domain.Post, error)
	GetPost(id string) (*domain.Post, error)
	GetPosts() ([]*domain.Post, error)
	AddComment(postID, parentID, content string) (*domain.Comment, error)
}
