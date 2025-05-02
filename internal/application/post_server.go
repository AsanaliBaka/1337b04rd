package application

import "1337b04rd/internal/domain"

type PostService struct {
	postRepo     domain.PostRepository
	imageStorage domain.ImageStorage
}

func NewPostService(postRepo domain.PostRepository, imageStorage domain.ImageStorage) *PostService {
	return &PostService{postRepo: postRepo, imageStorage: imageStorage}
}
