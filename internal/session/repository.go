package session

import (
	"context"
	"database/sql"
	"fmt"
)

type SessionRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (s *SessionRepo) AddSession(ctx context.Context, session Session) error {
	query := `INSERT INTO sessions (id, user_id, date_created, csrf_token)
	VALUES (?, ?, ?, ?);`

	_, err := s.db.ExecContext(ctx, query, session.id, session.userID, session.createdAt, session.csrfToken)
	if err != nil {
		return fmt.Errorf("adding session: %w", err)
	}

	return nil
}
