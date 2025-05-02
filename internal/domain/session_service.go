package domain

import (
	"context"
	"log"
	"sync"
	"time"
)

type Session struct {
	PostID    string
	ExpiresAt time.Time
}

type SessionService struct {
	sessions map[string]Session
	mu       sync.RWMutex
	postRepo PostRepository
}

func NewSessionService(postRepo PostRepository) *SessionService {
	s := &SessionService{
		sessions: make(map[string]Session),
		postRepo: postRepo,
	}
	go s.watchSessions()
	return s
}

func (s *SessionService) AddSession(postID string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[postID] = Session{PostID: postID, ExpiresAt: expiresAt}
}

func (s *SessionService) UpdateSession(postID string, newExpiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.sessions[postID]; exists {
		s.sessions[postID] = Session{PostID: postID, ExpiresAt: newExpiresAt}
	}
}

func (s *SessionService) watchSessions() {
	for {
		time.Sleep(time.Second)
		s.mu.Lock()
		for id, session := range s.sessions {
			if time.Now().After(session.ExpiresAt) {
				err := s.postRepo.UpdatePost(context.Background(), id, true, time.Now())
				if err != nil {
					log.Printf("ошибка при архивировании поста %s: %v", id, err)
				}
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}
