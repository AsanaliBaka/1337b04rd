// internal/adapters/right/db/postgres/postgres.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"1337b04rd/internal/domain"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) CreatePost(post *domain.Post) error {
	query := `
        INSERT INTO posts (id, title, content, image_url, user_name, user_avatar, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.Exec(
		query,
		post.ID,
		post.Title,
		post.Content,
		post.ImageURL,
		post.User.Name,
		post.User.AvatarURL,
		post.CreatedAt,
	)

	return err
}

func (r *PostgresRepository) GetPostByID(id string) (*domain.Post, error) {
	query := `
        SELECT id, title, content, image_url, user_name, user_avatar, created_at
        FROM posts
        WHERE id = $1
    `

	var post domain.Post
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.ImageURL,
		&post.User.Name,
		&post.User.AvatarURL,
		&post.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Загружаем комментарии для поста
	comments, err := r.getCommentsForPost(id)
	if err != nil {
		return nil, err
	}

	post.Comments = comments
	return &post, nil
}

func (r *PostgresRepository) GetPosts() ([]*domain.Post, error) {
	query := `
        SELECT id, title, content, image_url, user_name, user_avatar, created_at
        FROM posts
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.ImageURL,
			&post.User.Name,
			&post.User.AvatarURL,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}

		// Для списка постов не загружаем комментарии (ленивая загрузка)
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostgresRepository) CreateComment(comment *domain.Comment) error {
	query := `
        INSERT INTO comments (id, post_id, parent_id, content, user_name, user_avatar, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.Exec(
		query,
		comment.ID,
		comment.PostID,
		comment.ParentID,
		comment.Content,
		comment.User.Name,
		comment.User.AvatarURL,
		comment.CreatedAt,
	)

	return err
}

func (r *PostgresRepository) getCommentsForPost(postID string) ([]domain.Comment, error) {
	query := `
        SELECT id, post_id, parent_id, content, user_name, user_avatar, created_at
        FROM comments
        WHERE post_id = $1
        ORDER BY created_at ASC
    `

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Content,
			&comment.User.Name,
			&comment.User.AvatarURL,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *PostgresRepository) DeleteOldPosts(noCommentsThreshold, withCommentsThreshold time.Duration) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Удаляем посты без комментариев старше порога
	_, err = tx.ExecContext(ctx, `
        DELETE FROM posts
        WHERE id IN (
            SELECT p.id
            FROM posts p
            LEFT JOIN comments c ON p.id = c.post_id
            WHERE c.id IS NULL
            AND p.created_at < $1
        )
    `, time.Now().Add(-noCommentsThreshold))
	if err != nil {
		return err
	}

	// Удаляем посты с комментариями, где последний комментарий старше порога
	_, err = tx.ExecContext(ctx, `
        DELETE FROM posts
        WHERE id IN (
            SELECT p.id
            FROM posts p
            JOIN (
                SELECT post_id, MAX(created_at) as last_comment_time
                FROM comments
                GROUP BY post_id
            ) c ON p.id = c.post_id
            WHERE c.last_comment_time < $1
        )
    `, time.Now().Add(-withCommentsThreshold))
	if err != nil {
		return err
	}

	return tx.Commit()
}
