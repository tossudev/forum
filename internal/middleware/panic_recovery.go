package middleware

import (
	"fmt"
	"forum/internal/errs"
	"net/http"
)

// RecoverPanic checks if a panic occurred and calls errs.WriteError to write an HTTP response
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Deferred function will always run in the event of a panic
		defer func() {
			panicValue := recover()

			if panicValue != nil {
				w.Header().Set("Connection", "close")
				errs.WriteError(w, fmt.Errorf("%v", panicValue))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
