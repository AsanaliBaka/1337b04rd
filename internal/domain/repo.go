package domain

import "context"

type PostRepository interface {
	CreatePost(ctx context.Context, post *Post) error
	GetByIDPost(ctx context.Context, id string) (*Post, error)
	GetAllPost(ctx context.Context) ([]*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, id string) error
}

type CommentRepository interface {
	CreateComment(comment *Comment) error
	GetByPostIDComment(postID string) ([]*Comment, error)
}
