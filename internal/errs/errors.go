package errs

import (
	"context"
	"encoding/json"
	"errors"
	"forum/internal/pagination"
	"log/slog"
	"net/http"
)

var (
	ErrInvalidFiletype = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrDuplicate       = errors.New("duplicate entry")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrSearchNoResults = errors.New("no results found")

	// Errors used for determining which error page to render
	ErrNotFound         = errors.New("record not found") // 404
	ErrInvalidUserInput = errors.New("invalid input")    // 400

	// User facing errors (error message visible to users)
	ErrHTMLSanitizationFail = errors.New("Invalid HTML input.")
	ErrInvalidThreadTitle   = errors.New("Thread title must be 1-100 characters.")
	ErrInvalidThreadBody    = errors.New("Thread body must be 1-5000 characters")
	ErrInvalidUsername      = errors.New("Username must be 2-20 characters, start with a letter, and only contain letters, numbers, or underscores (_).")
	ErrInvalidEmail         = errors.New("Enter a valid email address, such as example@email.com.")
	ErrInvalidPassword      = errors.New("Password must be 8-128 characters.")
	ErrDuplicateUsername    = errors.New("Username already taken.")
	ErrDuplicateEmail       = errors.New("Email already taken.")
	ErrInvalidCredentials   = errors.New("Invalid credentials.")
	ErrInternal             = errors.New("Sorry, we're unable to process your request at the moment. Please try again later.")
)

// WriteError chooses the appropriate error response and writes to http.ResponseWrite
func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		slog.Info("user disconnected before response finished")
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "request timed out", http.StatusGatewayTimeout)
	case errors.Is(err, ErrNotFound):
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidUserInput),
		errors.Is(err, pagination.ErrPaginationParams),
		errors.As(err, new(*json.SyntaxError)),
		errors.As(err, new(*json.UnmarshalTypeError)),
		errors.Is(err, ErrDuplicate):
		http.Error(w, ErrInvalidUserInput.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrInvalidFiletype):
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	case errors.Is(err, ErrFileTooLarge):
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
	case errors.Is(err, ErrUnauthorized):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		slog.Error("internal server error:", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
