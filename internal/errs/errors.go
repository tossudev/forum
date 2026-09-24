package errs

import (
	"context"
	"errors"
	"forum/internal/pagination"
	"log/slog"
	"net/http"
)

var (
	ErrInvalidUserInput = errors.New("invalid input")
  ErrInvalidFiletype  = errors.New("invalid file type")
  ErrFileTooLarge     = errors.New("file too large")
	ErrNotFound         = errors.New("record not found")
	ErrDuplicate        = errors.New("duplicate entry")
	ErrUnauthorized     = errors.New("unauthorized")
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

  case errors.Is(err, ErrFileTooLarge):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, ErrUnauthorized):
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return

	default:
		slog.Error("internal server error:", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
