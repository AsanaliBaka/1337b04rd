package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"1337b04rd/internal/domain"
)

func (app *App) AddComment(ctx context.Context, postID string, comment *domain.Comment) error {
	_, err := app.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found from db: %w", err)
	}

	author, err := app.repo.GetUserByID(ctx, comment.Author)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	comment.AvatarLink = author.ImageURL

	app.Lock()
	defer app.Unlock()

	// Если нет активного таймера, значит пост "не активен"
	if _, ok := app.timers[postID]; !ok {
		return fmt.Errorf("post with ID %s is not active", postID)
	}

	// Останавливаем старый таймер, если был
	app.timers[postID].Stop()

	// Запускаем новый таймер
	app.timers[postID] = time.AfterFunc(15*time.Minute, func() {
		app.Lock()
		defer app.Unlock()

		if timer, ok := app.timers[postID]; ok {
			delete(app.timers, postID)
			go app.archivePost(context.Background(), postID)
			timer.Stop()
		}
	})

	err = app.repo.AddComment(ctx, postID, comment)
	if err != nil {
		return fmt.Errorf("failed to add comment in database: %w", err)
	}
	return nil
}

func (app *App) ReplyToComment(ctx context.Context, parentCommentID string, reply *domain.Comment) error {
	parentPost, err := app.findPostByCommentID(ctx, parentCommentID)
	if err != nil {
		return fmt.Errorf("parent comment not found: %w", err)
	}

	author, err := app.repo.GetUserByID(ctx, reply.Author)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	reply.AvatarLink = author.ImageURL

	// 2. Проверяем активность поста
	app.Lock()
	defer app.Unlock()

	if _, ok := app.timers[parentPost.ID]; !ok {
		return errors.New("cannot comment on archived post")
	}

	// 3. Устанавливаем родителя для ответа
	reply.ParentID = parentCommentID
	reply.CreatedAt = time.Now()

	// 4. Обновляем таймер активности поста (аналогично AddComment)
	app.resetPostTimer(parentPost.ID)

	// 5. Добавляем комментарий в хранилище
	if err := app.repo.AddComment(ctx, parentPost.ID, reply); err != nil {
		return fmt.Errorf("failed to add reply: %w", err)
	}

	return nil
}

// Вспомогательный метод для поиска поста по ID комментария
func (app *App) findPostByCommentID(ctx context.Context, commentID string) (*domain.Post, error) {
	posts, err := app.repo.GetPosts(ctx)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		if containsComment(post.Comments, commentID) {
			return &post, nil
		}
	}

	return nil, errors.New("comment not found in any post")
}

// Рекурсивный поиск комментария
func containsComment(comments []domain.Comment, targetID string) bool {
	for _, comment := range comments {
		if comment.ID == targetID {
			return true
		}
		if containsComment(comment.Replies, targetID) {
			return true
		}
	}
	return false
}

// Обновление таймера поста (вынесено в отдельный метод для reuse)
func (app *App) resetPostTimer(postID string) {
	if timer, exists := app.timers[postID]; exists {
		timer.Stop()
	}

	app.timers[postID] = time.AfterFunc(15*time.Minute, func() {
		app.Lock()
		defer app.Unlock()

		if timer, ok := app.timers[postID]; ok {
			delete(app.timers, postID)
			go app.archivePost(context.Background(), postID)
			timer.Stop()
		}
	})
}
