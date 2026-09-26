package thread

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"
	"forum/internal/session"
	//"forum/internal/entities/user"
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

	data := ThreadsPage{
		// TODO: use category name instead of ID
		Path:    fmt.Sprintf("⌂ Home / CategoryID %d", id),
		Threads: threads,
		Page:    pagination.Page,
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

	var input struct {
		Title      string `json:"title"`
		Body       string `json:"body"`
		CategoryID int    `json:"category_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
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
		Body:       input.Body,
		AuthorID:   authorID,
		CategoryID: input.CategoryID,
	}

	err = h.service.Create(ctx, &thread)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	//TODO: find out what front end needs this result to do/look like

	w.WriteHeader(http.StatusOK)
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

func (h *ThreadHandler) SearchThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
	}

	//TODO: how will search term actually be sent from front end??
	term := r.PathValue("search")

	threads, err := h.service.SearchThreads(ctx, term, pagination)
	if err != nil {
		errs.WriteError(w, err)
	}

	var username string
	//will need this after merge with main, where all responses to front end must contain username
	//user := user.GetUser(r)
	//if user != nil {
	//		username = user.Username
	//}

	type SearchResults struct {
		Threads  []Thread
		Page     int
		Username string
	}

	data := SearchResults{
		Threads:  threads,
		Page:     pagination.Page,
		Username: username,
	}

	//TODO: sub actual frontent page when it exists (if not "searchresults.html")
	h.renderer.RenderPage(w, "searchresults.html", data)
}
