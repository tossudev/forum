package category

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"forum/internal/errs"
	"forum/internal/pagination"
)

type CategoryHandler struct {
	service   *CategoryService
	validator *validator.Validate
}

func NewHandler(service *CategoryService, validator *validator.Validate) *CategoryHandler {
	return &CategoryHandler{service: service, validator: validator}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	_, err = h.service.GetAll(ctx, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	//TODO: change _ to category and send to front end
	_, err = h.service.GetByID(ctx, id)

	//TODO: find out what front end needs this result to do/look like

	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid request body", err))
		return
	}

	category := Category{
		Name: input.Name,
	}

	err = h.service.Create(ctx, &category)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	//TODO: find out what front end needs this result to do/look like

	w.WriteHeader(http.StatusOK)
}
