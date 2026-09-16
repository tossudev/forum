package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type SessionManager struct {
	repo           *SessionRepo
	cookieName     string
	idleExpiration time.Duration
}

type Session struct {
	id            string
	csrfToken     string
	userID        int
	createdAt     time.Time
	idleExpiresAt time.Time
}

// sessionContextKey is used as a context key for getting the session from context
// Having this separate type helps prevent potential key naming collisions
type sessionContextKey struct{}

func NewSessionManager(repo *SessionRepo, cookieName string, idleExpiration, absoluteExpiration time.Duration) *SessionManager {
	return &SessionManager{
		repo:           repo,
		cookieName:     cookieName,
		idleExpiration: idleExpiration,
	}
}

func (sm *SessionManager) NewSession(userID int) (*Session, error) {
	sessionID, err := generateToken(32) // 32 bytes or 256 bits of randomness
	if err != nil {
		return nil, fmt.Errorf("creating new session: %w", err)
	}

	csrfToken, err := generateToken(32) // 32 bytes or 256 bits of randomness
	if err != nil {
		return nil, fmt.Errorf("creating new session: %w", err)
	}

	now := time.Now()

	s := Session{
		id:            sessionID,
		csrfToken:     csrfToken,
		userID:        userID,
		createdAt:     now,
		idleExpiresAt: now.Add(sm.idleExpiration),
	}

	return &s, nil
}

func generateToken(length int) (string, error) {
	token := make([]byte, length)

	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}

func (sm *SessionManager) writeCookie(w http.ResponseWriter, session *Session) {
	cookie := &http.Cookie{
		Name:     sm.cookieName,
		Value:    session.id,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // TODO: Use when HTTPS is implemented
		Expires: session.idleExpiresAt,
		MaxAge:  int(sm.idleExpiration / time.Second),
	}

	http.SetCookie(w, cookie)
}

// Authenticate session middleware - check cookie for session ID, then validate it
func (sm *SessionManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var session *Session

		// Read session ID from cookie
		cookie, err := r.Cookie(sm.cookieName)
		if err == nil {
			sessionID := cookie.Value
			session, err = sm.repo.GetSessionByID(ctx, sessionID)
			if err != nil {
				slog.Error("failed to get session from repo", "err", err)
			}
		}

		// If the the session is expired, delete session
		if session != nil && session.isExpired() {
			sm.repo.DeleteSession(ctx, session.id)
			session = nil
		}

		// If the session is valid, update last idleExpiration
		if session != nil && !session.isExpired() {
			session.idleExpiresAt = time.Now().Add(sm.idleExpiration)
			sm.repo.UpdateExpiry(ctx, session.id)
		}

		// Attach session to context
		ctx = context.WithValue(ctx, sessionContextKey{}, session)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *Session) isExpired() bool {
	return s.idleExpiresAt.Before(time.Now())
}
