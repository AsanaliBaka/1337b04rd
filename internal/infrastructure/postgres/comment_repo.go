package postgres

import (
	"1337b04rd/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type commentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) domain.CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (c *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	_, err := c.db.Exec(ctx, qCreateComment,
		&comment.ID,
		&comment.Text,
		&comment.Author,
		&comment.PostId,
		&comment.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil

}
func (c *commentRepository) GetByPostIDComment(ctx context.Context, postID string) ([]*domain.Comment, error) {
	var items []*domain.Comment

	comments, err := c.db.Query(ctx, qGetByPostIDComment)

	if err != nil {
		return nil, err
	}

	defer comments.Close()

	for comments.Next() {
		item := new(domain.Comment)

		if err := comments.Scan(
			&item.ID,
			&item.Text,
			&item.Author,
			&item.PostId,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %w", err)

		}

		items = append(items, item)
	}

	return items, nil
}
