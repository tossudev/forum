package comment

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"
	"forum/internal/session"
	"forum/internal/validate"

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

	user := user.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || threadID < 1 {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	if err := paginationData.Validate(); err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	// input params for CollectPageData
	params := thread.ThreadPageParams{
		ThreadID:       threadID,
		PaginationData: paginationData,
		Username:       user.Username,
		Path:           "/threads/" + strconv.Itoa(threadID),
	}
	// to check if there is only spaces in comment body
	commentBody := validate.TrimEmpty(r.PostFormValue("body"))

	if commentBody == "" {
		pageData, _, err := h.threadService.CollectPageData(ctx, params)
		if err != nil {
			h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
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

	_, totalComments, err := h.service.Create(ctx, &comment)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUserInput) {
			pageData, _, err := h.threadService.CollectPageData(ctx, params)
			if err != nil {
				h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
				return
			}

			pageData.CommentBody = commentBody
			pageData.ErrorMessage = "Comment has invalid formatting"
			w.WriteHeader(http.StatusBadRequest)
			h.renderer.RenderPage(w, "thread.html", pageData)
			return
		}
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}

	lastPage := pagination.CalcTotalPages(totalComments, pagination.DefaultPageSize)

	// On successful submission, redirect to the last page of the thread
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID)+"?page="+strconv.Itoa(lastPage), http.StatusSeeOther)
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
