package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
