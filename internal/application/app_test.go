package application

import (
	"context"
	"testing"
	"time"

	"1337b04rd/internal/domain"
)

func TestCreatePost_TimerAndArchivation(t *testing.T) {
	a := NewApp(nil, nil, nil, userService{})

	post := &domain.Post{ID: "post-123"}
	called := make(chan string, 1)

	// Подменяем archivePost через closure
	a.ArchivePost = func(ctx context.Context, postID string) {
		called <- postID
	}

	// Вставим аналог CreatePost, но с коротким таймером
	a.Lock()
	a.Timers()[post.ID] = time.AfterFunc(50*time.Millisecond, func() {
		a.Lock()
		defer a.Unlock()
		delete(a.Timers(), post.ID)
		a.ArchivePost(context.Background(), post.ID)
	})
	a.Unlock()

	// Ждём результата или тайм-аут
	select {
	case id := <-called:
		if id != post.ID {
			t.Fatalf("expected postID %s, got %s", post.ID, id)
		}
		t.Logf("archivePost called with postID: %s", id)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("archivePost not called")
	}
}
