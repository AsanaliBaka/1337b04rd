package right

import (
	"context"

	"1337b04rd/internal/domain"
)

type DbPort interface {
	PostRepository
	ArchiveRepository
	CommentRepository
	UserRepository
	SessionRepository
}

type PostRepository interface {
	GetPosts(ctx context.Context) ([]domain.Post, error)
	ListCatalog(ctx context.Context) ([]*domain.PostSummary, error)
	GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) error
}

type ArchiveRepository interface {
	ListArchiveCatalog(ctx context.Context) ([]*domain.PostSummary, error)
	GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error)
	ArchivePostByID(ctx context.Context, id string) (*domain.Post, error)
}
type CommentRepository interface {
	AddComment(ctx context.Context, PostId string, comment *domain.Comment) error
	ReplyToComment(ctx context.Context, PostID string, UserID string, comment *domain.Comment) error
	// ReplyToComment(ctx context.Context, userID string, parentCommentID string, text string) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	GetMaxCharacterID(ctx context.Context) (int, error)
}

type SessionRepository interface {
	GetSession(ctx context.Context, sessionID string) (*domain.Session, error)
	SaveSession(ctx context.Context, session *domain.Session) error
}
