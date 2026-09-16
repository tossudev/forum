package session

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/errs"
)

type SessionRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) AddSession(ctx context.Context, session Session) error {
	query := `INSERT INTO sessions (id, user_id, csrf_token, date_created, expires_at)
	VALUES (?, ?, ?, ?);`

	_, err := r.db.ExecContext(ctx, query, session.id, session.userID, session.csrfToken, session.createdAt, session.idleExpiresAt)
	if err != nil {
		return fmt.Errorf("adding session: %w", err)
	}

	return nil
}

func (r *SessionRepo) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	query := `SELECT (id, user_id, csrf_token, date_created, expires_at)
	FROM sessions WHERE id = ?;`

	s := &Session{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&s.id, &s.userID, &s.csrfToken, &s.createdAt, &s.idleExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: invalid session ID", errs.ErrsUnauthorized)
		}
		return nil, fmt.Errorf("GetSessionByID: scanning row: %w", err)
	}

	return s, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?;`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}

	return nil
}

func (r *SessionRepo) UpdateExpiry(ctx context.Context, id string) error {
	query := `UPDATE sessions SET expires_at WHERE id = ?;`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("refreshing session expiry: %w", err)
	}

	return nil
}
