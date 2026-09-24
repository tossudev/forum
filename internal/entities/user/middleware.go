package user

import (
	"context"
	"forum/internal/session"
	"log/slog"
	"net/http"
)

type userContextKey struct{}

// UserMiddleware is a middleware that attaches the current *User to the request context if the session is authenticated (user is logged in)
func (h *UserHandler) GetUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := session.GetSession(r)

		var user *User
		var err error

		if session != nil {
			user, err = h.service.repo.GetUserByID(ctx, session.UserID())
			if err != nil {
				// TODO: redirect to error page according to error type
				return
			}
		}
		// Attach username to context
		ctx = context.WithValue(ctx, userContextKey{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetUser is a helper that allows handlers to access a *User stored in the request context.
// If the user is not logged in, *User will be nil.
func GetUser(r *http.Request) *User {
	user, ok := r.Context().Value(userContextKey{}).(*User)
	if !ok {
		slog.Error("user not found in request context", "path", r.URL.Path, "method", r.Method)
	}
	return user
}
