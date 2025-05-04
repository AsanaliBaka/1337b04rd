package postgres

import (
	"context"
	"fmt"

	"1337b04rd/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) domain.SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (s *sessionRepository) CreateSession(ctx context.Context, session *domain.UserRef) error {
	_, err := s.db.Exec(ctx, qCreateSession,
		session.SessionID,
		session.AvatarURL,
		session.Name)
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}
	return nil
}

func (s *sessionRepository) GetByIDSession(ctx context.Context, sessionID string) (*domain.UserRef, error) {
	var userSession domain.UserRef
	err := s.db.QueryRow(ctx, qGetByIDSession).Scan(
		&userSession.SessionID,
		&userSession.AvatarURL,
		&userSession.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	return &userSession, nil
}

func (s *sessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := s.db.Exec(ctx, qDeleteSession, sessionID)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}
	return nil
}
