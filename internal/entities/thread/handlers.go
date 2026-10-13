package thread

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"forum/internal/entities/user"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"
	"forum/internal/session"
)

type ThreadHandler struct {
	service   *ThreadService
	validator *validator.Validate
	renderer  *render.Renderer
}

func NewHandler(service *ThreadService, validator *validator.Validate, renderer *render.Renderer) *ThreadHandler {
	return &ThreadHandler{
		service:   service,
		validator: validator,
		renderer:  renderer,
	}
}

func (h *ThreadHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid category id", errs.ErrInvalidUserInput))
		return
	}

	threads, err := h.service.GetByCategory(ctx, id, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	data := ThreadsPage{
		// TODO: use category name instead of ID
		Path:     fmt.Sprintf("⌂ Home / Category: %d", id),
		Threads:  threads,
		Page:     pagination.Page,
		Username: username,
	}

	h.renderer.RenderPage(w, "threads.html", data)
}

func (h *ThreadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	//TODO: change _ to thread and send to front end
	_, err = h.service.GetByID(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ThreadHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	type submission struct {
		Title      string
		Body       string
		CategoryID int
	}

	// TODO: FIX IDIOT
	input := submission{
		Title:      r.PostFormValue("title"),
		Body:       r.PostFormValue("body"),
		CategoryID: 1,
	}

	// TODO: sanitizing SHOULD be done only upon displaying it
	// so instead of doing it here, we should store the raw input in db
	// and sanitize the body when we want to display it
	sanitizedBody, err := h.renderer.SanitizeHTML(input.Body)
	if err != nil {
		// Could be status 422: Unprocessable Content
		errs.WriteError(w, fmt.Errorf("%w: invalid request body", errs.ErrInvalidUserInput))
		return
	}

	session := session.GetSession(r)
	if session == nil {
		errs.WriteError(w, fmt.Errorf("%w: not authenticated", errs.ErrUnauthorized))
		return
	}
	authorID := session.UserID()

	thread := Thread{
		Title:      input.Title,
		Body:       sanitizedBody,
		AuthorID:   authorID,
		CategoryID: input.CategoryID,
	}

	id, err := h.service.Create(ctx, &thread)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	// Redirect user to newly created thread page
	http.Redirect(w, r, fmt.Sprintf("/threads/%d", id), http.StatusSeeOther)
}

func (h *ThreadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session := session.GetSession(r)
	if session == nil {
		errs.WriteError(w, fmt.Errorf("%w: not authenticated", errs.ErrUnauthorized))
		return
	}
	userID := session.UserID()

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	err = h.service.Delete(ctx, threadID, userID)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
