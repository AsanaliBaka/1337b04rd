package postgres

import (
	"context"
	"fmt"
	"time"

	"1337b04rd/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) domain.PostRepository {
	return &postRepo{
		db: db,
	}
}

func (p *postRepo) CreatePost(ctx context.Context, post *domain.Post) error {
	_, err := p.db.Exec(ctx, qCreatePost,
		post.ID,
		post.Title,
		post.Content,
		post.ImageURL,
		post.Author.SessionID,
		post.CreatedAt,
		post.UpdatedAt,
		post.IsArchived,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *postRepo) GetByIDPost(ctx context.Context, id string) (*domain.Post, error) {
	var item domain.Post

	err := p.db.QueryRow(ctx, qGetByIDPost, id).Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.ImageURL,
		&item.Author,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.IsArchived,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (p *postRepo) GetAllPost(ctx context.Context) ([]*domain.Post, error) {
	rows, err := p.db.Query(ctx, qGetAllPost)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.Post

	for rows.Next() {
		var item domain.Post

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.ImageURL,
			&item.Author.SessionID,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.IsArchived,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (p *postRepo) UpdatePost(ctx context.Context, postId string, isArchive bool, upadtedAt time.Time) error {
	_, err := p.db.Exec(ctx, qUpdatePost, upadtedAt, isArchive, postId)
	if err != nil {
		return fmt.Errorf("failed to update post %s: %w", postId, err)
	}
	return nil
}

func (p *postRepo) DeletePost(ctx context.Context, id string) error {
	_, err := p.db.Exec(ctx, qDeletePost, id)
	if err != nil {
		return err
	}

	return nil
}
