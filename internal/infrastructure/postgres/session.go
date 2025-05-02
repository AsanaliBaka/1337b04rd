package postgres

import (
	. "1337b04rd/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSession(db *pgxpool.Pool) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}
func (s *sessionRepository) CreateSession(ctx context.Context, session *UserRef) error {
	_, err := s.db.Exec(ctx, qCreateSession,
		session.SessionID,
		session.AvatarURL,
		session.Name)

	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}
	return nil
}
func (s *sessionRepository) GetByIDSession(ctx context.Context, sessionID string) (*UserRef, error) {
	var userSession UserRef
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
