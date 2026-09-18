package thread

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/session"
)

type ThreadHandler struct {
	service   *ThreadService
	validator *validator.Validate
}

func NewHandler(service *ThreadService, validator *validator.Validate) *ThreadHandler {
	return &ThreadHandler{service: service, validator: validator}
}

func (h *ThreadHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	idString := r.URL.Query().Get("id")
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid category id", errs.ErrInvalidUserInput))
		return
	}

	//TODO: change _ to threads and send to front end
	_, err = h.service.GetByCategory(ctx, id, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	session, ok := session.GetSession(ctx)
	if session == nil || !ok {
		errs.WriteError(w, fmt.Errorf("%w: not authenticated", errs.ErrsUnauthorized))
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
