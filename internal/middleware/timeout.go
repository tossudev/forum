package middleware

import (
	"context"
	"net/http"
	"time"
)

// Timeout returns a middleware that passes a request with a timed context to the next handler
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	// Return a middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			// Pass control to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
