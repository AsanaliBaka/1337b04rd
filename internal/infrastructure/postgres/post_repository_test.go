package postgres_test

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/infrastructure/postgres"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepo(t *testing.T) {
	ctx := context.Background()
	pool := setupDb(t)
	repo := postgres.NewPostRepo(pool)

	testAuthor := domain.UserRef{
		SessionID: "test_session_123",
		AvatarURL: "http://example.com/avatar.jpg",
		Name:      "Test User",
	}

	// Вставка пользователя в таблицу users
	_, err := pool.Exec(ctx, `
		INSERT INTO users(session_id, avatar_url, name)
		VALUES ($1, $2, $3)
	`, testAuthor.SessionID, testAuthor.AvatarURL, testAuthor.Name)
	require.NoError(t, err, "Failed to insert user for test")

	testPost := &domain.Post{
		ID:        "test_post_1",
		Title:     "Test Post",
		Content:   "This is a test post content",
		ImageURL:  "http://example.com/image.jpg",
		Author:    testAuthor,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Create and Get", func(t *testing.T) {
		err := repo.CreatePost(ctx, testPost)
		require.NoError(t, err, "CreatePost should not return error")

		fetch, err := repo.GetByIDPost(ctx, testPost.ID)
		require.NoError(t, err, "GetByIDPost should not return error")

		assert.Equal(t, testPost.ID, fetch.ID)
		assert.Equal(t, testPost.Title, fetch.Title)
		assert.Equal(t, testPost.Author.SessionID, fetch.Author.SessionID)
	})

	t.Run("GetAll", func(t *testing.T) {
		secondPost := &domain.Post{
			ID:        "test_post_2",
			Title:     "Second Post",
			Content:   "Second post content",
			ImageURL:  "http://example.com/image2.jpg",
			Author:    testAuthor,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		require.NoError(t, repo.CreatePost(ctx, secondPost))

		all, err := repo.GetAllPost(ctx)
		require.NoError(t, err)

		assert.Equal(t, 2, len(all), "Should return 2 posts")
	})

	t.Run("Update", func(t *testing.T) {
		newTime := time.Now().Add(1 * time.Hour)

		err := repo.UpdatePost(ctx, testPost.ID, true, newTime)
		require.NoError(t, err)

		updatedPost, err := repo.GetByIDPost(ctx, testPost.ID)
		require.NoError(t, err)

		assert.True(t, updatedPost.IsArchived, "Post should be archived")
		assert.WithinDuration(t, newTime, updatedPost.UpdatedAt, time.Second)
	})
}

func setupDb(t *testing.T) *pgxpool.Pool {

	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), conn)
	require.NoError(t, err, "Failed to connect to test database")

	err = pool.Ping(context.Background())
	require.NoError(t, err, "Failed to ping database")

	var whoami string
	err = pool.QueryRow(context.Background(), "SELECT current_user").Scan(&whoami)
	require.NoError(t, err)
	t.Logf("Connected as: %s", whoami)

	// Очистка и создание таблиц
	_, err = pool.Exec(context.Background(), `
		DROP TABLE IF EXISTS posts;
		DROP TABLE IF EXISTS users;

		CREATE TABLE users(
			session_id VARCHAR(36) PRIMARY KEY,
			avatar_url TEXT,
			name TEXT
		);

		CREATE TABLE posts(
			id VARCHAR(50) PRIMARY KEY, 
			title TEXT NOT NULL, 
			content TEXT NOT NULL,
			image_url TEXT NOT NULL,
			author_id VARCHAR(36) NOT NULL REFERENCES users(session_id),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE,
			isarchived BOOLEAN DEFAULT FALSE
		);
	`)
	require.NoError(t, err, "Failed to create test tables")

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
