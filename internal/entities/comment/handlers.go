package comment

import (
	"errors"
	"fmt"
	"forum/internal/errs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/entities/user"
	"forum/internal/render"
	"forum/internal/session"

	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	service   *CommentService
	validator *validator.Validate
	renderer  *render.Renderer
}

func NewHandler(service *CommentService, validator *validator.Validate, renderer *render.Renderer) *CommentHandler {
	return &CommentHandler{service: service, validator: validator, renderer: renderer}
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

	//ErrorPage is used to relay values to the frontend in error pages
	ep := ErrorPage{
		Username: user.Username,
		Message:  "Comment has to include characters",
	}

	// to check if there is only spaces in comment body
	commentBody := strings.TrimSpace(r.PostFormValue("body"))
	if commentBody == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "error-400.html", ep)
		return
	}
	comment := Comment{
		Body:     commentBody,
		ThreadID: threadID,
		AuthorID: user.ID,
	}
	//call to service layer
	_, err = h.service.Create(ctx, &comment)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUserInput) {
			w.WriteHeader(http.StatusBadRequest)
			ep.Message = "comment has invalid formatting"
			h.renderer.RenderPage(w, "error-400.html", ep)
			return
		}
		errs.WriteError(w, err)
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
