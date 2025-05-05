package application

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/ports/left"
)

type PostService struct {
	postRepo     left.PostRepository
	avatarRepo   left.AvatarProvider
	imageStorage left.ImageStorage
}

func NewPostService(pr left.PostRepository, ar left.AvatarProvider, is left.ImageStorage) left.PostService {
	return &PostService{
		postRepo:     pr,
		avatarRepo:   ar,
		imageStorage: is,
	}
}

func (s *PostService) CreatePost(title, content string, image []byte) (*domain.Post, error) {
	return nil, nil
}

func (s *PostService) GetPost(id string) (*domain.Post, error) {
	return nil, nil
}

func (s *PostService) GetPosts() ([]*domain.Post, error) {
	return nil, nil
}

func (s *PostService) AddComment(postID, parentID, content string) (*domain.Comment, error) {
	return nil, nil
}
