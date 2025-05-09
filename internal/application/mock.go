package application

import (
	"context"
	"errors"
	"time"

	"1337b04rd/internal/domain"
)

type MockPostRepository struct{}

func (m *MockPostRepository) ListCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	// Возвращаем тестовые данные
	return []*domain.PostSummary{
		{
			ID:    "1",
			Title: "Post 1",
		},
		{
			ID:    "2",
			Title: "Post 2",
		},
	}, nil
}

func (m *MockPostRepository) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	if id == "valid-post" {
		return &domain.Post{
			ID:        "valid-post",
			Title:     "Valid Post",
			Content:   "This is a valid post.",
			CreatedAt: time.Now(),
		}, nil
	}
	return nil, errors.New("post not found")
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post *domain.Post) error {
	// Предполагаем успешное создание поста
	return nil
}

type MockArchiveRepository struct{}

func (m *MockArchiveRepository) ListArchiveCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	return []*domain.PostSummary{
		{
			ID:    "3",
			Title: "Archived Post 1",
		},
	}, nil
}

func (m *MockArchiveRepository) GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error) {
	if id == "archived-post" {
		return &domain.Post{
			ID:        "archived-post",
			Title:     "Archived Post",
			Content:   "This is an archived post.",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		}, nil
	}
	return nil, errors.New("archived post not found")
}

func (m *MockArchiveRepository) ArchivePostByID(ctx context.Context, id string) (*domain.Post, error) {
	// Возвращаем nil, чтобы имитировать успешную архивацию
	return nil, nil
}

// Mock для CommentRepository
type MockCommentRepository struct{}

func (m *MockCommentRepository) AddComment(ctx context.Context, PostId string, comment *domain.Comment) error {
	if PostId == "valid-post" {
		return nil // Успешное добавление комментария
	}
	return errors.New("post not found")
}

func (m *MockCommentRepository) ReplyToComment(ctx context.Context, PostID string, UserID string, comment *domain.Comment) error {
	return nil // Успешный ответ на комментарий
}

// Mock для UserRepository
type MockUserRepository struct{}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return nil // Имитация успешного создания пользователя
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "user-123" {
		return &domain.User{
			ID:       "user-123",
			Username: "testuser",
		}, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetMaxCharacterID(ctx context.Context) (int, error) {
	return 100, nil // Пример возврата максимального ID
}

// Mock для SessionRepository
type MockSessionRepository struct{}

func (m *MockSessionRepository) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	if sessionID == "valid-session" {
		return &domain.Session{
			ID:        "valid-session",
			UserID:    "user-123",
			CreatedAt: time.Now().Add(-5 * time.Minute),
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}, nil
	}
	return nil, errors.New("session not found")
}

func (m *MockSessionRepository) SaveSession(ctx context.Context, session *domain.Session) error {
	// Имитация успешного сохранения сессии
	return nil
}

func (m *MockPostRepository) GetAllPosts(ctx context.Context) ([]domain.Post, error) {
	// Возвращаем тестовые данные
	return []domain.Post{
		{
			ID:        "1",
			Title:     "Post 1",
			Content:   "This is the first post.",
			CreatedAt: time.Now(),
		},
		{
			ID:        "2",
			Title:     "Post 2",
			Content:   "This is the second post.",
			CreatedAt: time.Now(),
		},
	}, nil
}

// Теперь создадим Mock для всего DbPort
type MockDbPort struct {
	MockPostRepository
	MockArchiveRepository
	MockCommentRepository
	MockUserRepository
	MockSessionRepository
}
