package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

type SessionManager struct {
	repo               *SessionRepo
	cookieName         string
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
}

type Session struct {
	id                string
	csrfToken         string
	userID            int
	createdAt         time.Time
	idleExpiresAt     time.Time
	absoluteExpiresAt time.Time
}

func NewSessionManager(repo *SessionRepo, cookieName string, idleExpiration, absoluteExpiration time.Duration) *SessionManager {
	return &SessionManager{
		repo:               repo,
		cookieName:         cookieName,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
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
		id:                sessionID,
		csrfToken:         csrfToken,
		userID:            userID,
		createdAt:         now,
		idleExpiresAt:     now.Add(sm.idleExpiration),
		absoluteExpiresAt: now.Add(sm.absoluteExpiration),
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
