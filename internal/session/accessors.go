package session

import "context"

// GetSession allows handlers to access a *Session stored in the context
// under the key, sessionContextKey{}. Returns false if *Session is not found.
func GetSession(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(*Session) // type assertion
	return session, ok
}

// UserID is a getter a Session's user ID
func (s *Session) UserID() int {
	return s.userID
}

// CSRF is a getter for a Session's csrf token
func (s *Session) CSRF() string {
	return s.csrfToken
}
