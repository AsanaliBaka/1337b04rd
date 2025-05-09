package left

import (
	"context"

	"1337b04rd/internal/domain"
)

type APIPort interface {
	PostQueryPort
	PostCommandPort
	SessionPort
}

type PostQueryPort interface {
	GetCatalog(ctx context.Context) ([]*domain.PostSummary, error)
	GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	GetArchiveList(ctx context.Context) ([]*domain.PostSummary, error)
	GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error)
}

type PostCommandPort interface {
	AddComment(ctx context.Context, postID string, comment *domain.Comment) error
	ReplyToComment(ctx context.Context, parentCommentID string, reply *domain.Comment) error
	CreatePost(ctx context.Context, post *domain.Post) error
}

type SessionPort interface {
	GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error)
	CreateSession(ctx context.Context) (*domain.Session, error)
}
