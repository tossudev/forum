package session

import (
	"log/slog"
	"net/http"
)

// GetSession allows handlers to access a *Session stored in the request context
// under the key, sessionContextKey{}.
func GetSession(r *http.Request) *Session {
	session, ok := r.Context().Value(sessionContextKey{}).(*Session) // type assertion
	if !ok {
		slog.Error("session not found in request context", "path", r.URL.Path, "method", r.Method)
	}
	return session
}

// UserID is a getter a Session's user ID
func (s *Session) UserID() int {
	return s.userID
}

// CSRF is a getter for a Session's csrf token
func (s *Session) CSRF() string {
	return s.csrfToken
}
