package errs

import (
	"context"
	"errors"
	"forum/internal/pagination"
	"log"
	"net/http"
)

var ErrInvalidUserInput = errors.New("invalid input")
var ErrNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("duplicate entry")

// WriteError chooses the appropriate error response and writes to http.ResponseWrite
func WriteError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		log.Println("user disconnected before response finished")
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, "request timed out", http.StatusGatewayTimeout)
		return
	}

	if errors.Is(err, ErrNotFound) {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if errors.Is(err, ErrInvalidUserInput) {
		http.Error(w, err.Error(), http.StatusBadRequest) // Show detailed error so user can fix input
		return
	}

	if errors.Is(err, pagination.ErrPaginationParams) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, ErrDuplicate) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("internal server error:", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
