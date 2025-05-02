package domain

import (
	"context"
	"io"
	"time"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *Post) error
	GetByIDPost(ctx context.Context, id string) (*Post, error)
	GetAllPost(ctx context.Context) ([]*Post, error)
	UpdatePost(ctx context.Context, postId string, isArchive bool, upadtedAt time.Time) error
	DeletePost(ctx context.Context, id string) error
}

type CommentRepository interface {
	CreateComment(comment *Comment) error
	GetByPostIDComment(postID string) ([]*Comment, error)
}

type ImageStorage interface {
	CreateImage(ctx context.Context, imageName string, imageData io.Reader, size int64) (string, error)
	GetImageURL(ctx context.Context, imageName string) (string, error)
	DeleteImage(ctx context.Context, imageName string) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *UserRef) error
	GetByIDSession(ctx context.Context, sessionID string) (*UserRef, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
