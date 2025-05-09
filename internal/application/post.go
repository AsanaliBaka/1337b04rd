package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"1337b04rd/internal/domain"

	"github.com/jackc/pgx"
)

func (app *App) CreatePost(ctx context.Context, post *domain.Post) error {
	app.Lock()
	defer app.Unlock()

	if err := app.repo.CreatePost(ctx, post); err != nil {
		return err
	}

	if timer, ok := app.timers[post.ID]; ok {
		timer.Stop()
	}
	app.timers[post.ID] = time.AfterFunc(10*time.Minute, func() {
		app.Lock()
		defer app.Unlock()

		if timer, ok := app.timers[post.ID]; ok {
			delete(app.timers, post.ID)
			go app.archivePost(context.Background(), post.ID)
			timer.Stop()
		}
	})

	return nil
}

func (app *App) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	post, err := app.repo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	author, err := app.repo.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	post.UserAvatar = author.ImageURL

	return post, nil
}

func (app *App) GetCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	catalog, err := app.repo.ListCatalog(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return catalog, nil
}

func (app *App) GetArchiveList(ctx context.Context) ([]*domain.PostSummary, error) {
	catalog, err := app.repo.ListArchiveCatalog(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return catalog, nil
}

func (app *App) GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error) {
	post, err := app.repo.GetArchivedPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	author, err := app.repo.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	post.UserAvatar = author.ImageURL

	return post, nil
}

func (app *App) archivePost(ctx context.Context, id string) (*domain.Post, error) {
	post, err := app.repo.ArchivePostByID(ctx, id)
	if err != nil {
		fmt.Printf("Failed to archive post %s: %v\n", id, err)
		return nil, err
	}
	fmt.Printf("Successfully archived post %s\n", id)
	return post, nil
}
