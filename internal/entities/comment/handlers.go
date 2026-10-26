package comment

import (
	"fmt"
	"forum/internal/errs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/entities/user"
	"forum/internal/session"

	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	service   *CommentService
	validator *validator.Validate
}

func NewHandler(service *CommentService, validator *validator.Validate) *CommentHandler {
	return &CommentHandler{service: service, validator: validator}
}

// TODO: user friendly error handling
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := user.GetUser(r)

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || threadID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("comment creation fail", "err", err)
		return
	}

	// TODO: this should be implemented for comments for better UX:
	// Parse form before handling user not logged in error so that form data can be saved
	// and reinserted into the page templates when displaying errors.
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("comment creation fail", "err", err)
		return
	}

	// User not logged in
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	comment := Comment{
		Body:     strings.TrimSpace(r.PostFormValue("body")),
		ThreadID: threadID,
		AuthorID: user.ID,
	}

	_, err = h.service.Create(ctx, &comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("comment creation fail", "err", err)
		return
	}

	// On successful submission, reload the page
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusSeeOther)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	if sess == nil {
		errs.WriteError(w, fmt.Errorf("%w: login required", errs.ErrUnauthorized))
		return
	}

	userID := sess.UserID()

	idString := r.PathValue("id")
	commentID, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid id", errs.ErrInvalidUserInput))
		return
	}

	err = h.service.Delete(r.Context(), commentID, userID)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
