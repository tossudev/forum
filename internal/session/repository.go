package session

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/errs"
	"time"
)

type SessionRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) addSession(ctx context.Context, session *Session) error {
	query := `INSERT INTO sessions (id, user_id, session_hash, csrf_token, date_created, expires_at)
	VALUES (?, ?, ?, ?, ?, ?);`

	createdAtStr := session.createdAt.Format(config.TimeFormat)
	expiresAtStr := session.idleExpiresAt.Format(config.TimeFormat)

	_, err := r.db.ExecContext(ctx, query, session.id, session.userID, session.sessionHash, session.csrfToken, createdAtStr, expiresAtStr)
	if err != nil {
		return fmt.Errorf("adding session: %w", err)
	}

	return nil
}

func (r *SessionRepo) getSessionByID(ctx context.Context, id string) (*Session, error) {
	query := `SELECT id, user_id, session_hash, csrf_token, date_created, expires_at
	FROM sessions WHERE id = ?;`

	var s Session
	var createdAtStr string
	var expiresAtStr string
	var parseErr error

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&s.id, &s.userID, &s.sessionHash, &s.csrfToken, &createdAtStr, &expiresAtStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: invalid session ID", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("getSessionByID: scanning row: %w", err)
	}

	s.createdAt, parseErr = time.Parse(config.TimeFormat, createdAtStr)
	if parseErr != nil {
		return nil, fmt.Errorf("getSessionByID: created at time parsing failure: %w", parseErr)
	}

	s.idleExpiresAt, parseErr = time.Parse(config.TimeFormat, expiresAtStr)
	if parseErr != nil {
		return nil, fmt.Errorf("getSessionByID: expires at time parsing failure: %w", parseErr)
	}

	return &s, nil
}

func (r *SessionRepo) deleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?;`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}

	return nil
}

func (r *SessionRepo) updateExpiry(ctx context.Context, session *Session) error {
	query := `UPDATE sessions SET expires_at = ? WHERE id = ?;`

	expiresAtStr := session.idleExpiresAt.Format(config.TimeFormat)
	if _, err := r.db.ExecContext(ctx, query, expiresAtStr, session.id); err != nil {
		return fmt.Errorf("refreshing session expiry: %w", err)
	}

	return nil
}
