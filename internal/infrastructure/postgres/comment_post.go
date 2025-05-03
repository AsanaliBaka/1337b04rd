package postgres

import (
	"context"
	"fmt"

	. "1337b04rd/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type commentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) CommentRepository {
	return &commentRepo{
		db: db,
	}
}

const (
	qCreateComment = `INSERT INTO comments 
		(id, text, author_session_id, image_url, post_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	qGetCommentsByPost = `SELECT id, text, author_session_id, image_url, created_at 
		FROM comments WHERE post_id = $1 ORDER BY created_at DESC`
)

func (c *commentRepo) CreateComment(ctx context.Context, comment *Comment) error {
	_, err := c.db.Exec(ctx, qCreateComment,
		comment.ID,
		comment.Text,
		comment.Author.SessionID,
		comment.ImageURL,
		comment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

// какая та хуйня с этим кодом, переделаю еще я по другому понял сперва а потом уже догнал
func (c *commentRepo) GetByPostIDComment(ctx context.Context, postID string) ([]*Comment, error) {
	var comments []*Comment

	rows, err := c.db.Query(ctx, qGetCommentsByPost, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.Author.SessionID,
			&comment.ImageURL,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}
