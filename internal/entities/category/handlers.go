package category

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"
)

type CategoryHandler struct {
	service   *CategoryService
	validator *validator.Validate
	renderer  *render.Renderer
}

func NewHandler(service *CategoryService, validator *validator.Validate, renderer *render.Renderer) *CategoryHandler {
	return &CategoryHandler{
		service:   service,
		validator: validator,
		renderer:  renderer,
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	categories, err := h.service.GetAll(ctx, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	data := CategoriesPage{
		Path:       "⌂ Home",
		Categories: categories,
		Page:       pagination.Page,
	}

	h.renderer.RenderPage(w, "landing.html", data)
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
