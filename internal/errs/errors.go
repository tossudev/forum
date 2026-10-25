package errs

import (
	"context"
	"errors"
	"forum/internal/pagination"
	"log/slog"
	"net/http"
)

var (
	ErrInvalidFiletype    = errors.New("invalid file type")
	ErrFileTooLarge       = errors.New("file too large")
	ErrDuplicate          = errors.New("duplicate entry")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrSearchNoResults    = errors.New("no results found")
	ErrInvalidThreadTitle = errors.New("invalid thread title")
	ErrInvalidThreadBody  = errors.New("invalid thread body")

	// Errors used for determining which error page to render
	ErrNotFound         = errors.New("record not found") // 404
	ErrInvalidUserInput = errors.New("invalid input")    // 400

	// User facing errors (error message visible to users)
	ErrInvalidUsername    = errors.New("Username must be 2-20 characters, start with a letter, and only contain letters, numbers, or underscores (_).")
	ErrInvalidEmail       = errors.New("Enter a valid email address, such as example@email.com.")
	ErrInvalidPassword    = errors.New("Password must be 8-128 characters.")
	ErrDuplicateUsername  = errors.New("Username already taken.")
	ErrDuplicateEmail     = errors.New("Email already taken.")
	ErrInvalidCredentials = errors.New("Invalid credentials.")
	ErrInternal           = errors.New("Sorry, we're unable to process your request at the moment. Please try again later.")
)

// WriteError chooses the appropriate error response and writes to http.ResponseWrite
func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		slog.Info("user disconnected before response finished")
		return

	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "request timed out", http.StatusGatewayTimeout)
		return

	case errors.Is(err, ErrNotFound):
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
		return

	case errors.Is(err, ErrInvalidUserInput), errors.Is(err, pagination.ErrPaginationParams), errors.Is(err, ErrDuplicate):
		http.Error(w, err.Error(), http.StatusBadRequest) // Show detailed error so user can fix input
		return

	case errors.Is(err, ErrInvalidFiletype):
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return

	case errors.Is(err, ErrFileTooLarge):
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		return

	case errors.Is(err, ErrUnauthorized):
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return

	default:
		slog.Error("internal server error:", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
