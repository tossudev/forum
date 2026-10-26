package comment

import (
	"errors"
	"fmt"
	"forum/internal/errs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/pagination"
	"forum/internal/render"
	"forum/internal/session"

	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	service       *CommentService
	validator     *validator.Validate
	renderer      *render.Renderer
	threadService *thread.ThreadService
}

func NewHandler(service *CommentService, validator *validator.Validate, renderer *render.Renderer, threadService *thread.ThreadService) *CommentHandler {
	return &CommentHandler{service: service, validator: validator, renderer: renderer, threadService: threadService}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	params := thread.Params{}

	user := user.GetUser(r)
	//user not logged in
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params.Username = user.Username

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || threadID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("comment creation fail", "err", err)
		return
	}
	params.ThreadID = threadID
	params.Path = "/threads/" + strconv.Itoa(threadID)

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", thread.ThreadCommentsPage{Username: params.Username})
		return
	}
	params.PaginationData = paginationData

	if err := paginationData.Validate(); err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", thread.ThreadCommentsPage{Username: params.Username})
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

	// to check if there is only spaces in comment body
	commentBody := strings.TrimSpace(r.PostFormValue("body"))
	if commentBody == "" {
		pageData, err := h.threadService.CollectPageData(ctx, params)
		if err != nil {
			errs.WriteError(w, err)
			return
		}
		pageData.CommentBody = commentBody
		pageData.ErrorMessage = "Comment has to include characters"
		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "thread.html", pageData)
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
			pageData, err := h.threadService.CollectPageData(ctx, params)
			if err != nil {
				errs.WriteError(w, err)
				return
			}
			pageData.CommentBody = commentBody
			pageData.ErrorMessage = "Comment has invalid formatting"
			w.WriteHeader(http.StatusBadRequest)
			h.renderer.RenderPage(w, "thread.html", pageData)
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
