package search

import (
	"fmt"
	"net/http"

	"forum/internal/entities/user"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"

	"github.com/go-playground/validator/v10"
)

type SearchHandler struct {
	service   *SearchService
	renderer  *render.Renderer
	validator *validator.Validate
}

func NewHandler(service *SearchService, renderer *render.Renderer, validator *validator.Validate) *SearchHandler {
	return &SearchHandler{service: service,
		renderer:  renderer,
		validator: validator}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	query := r.URL.Query()
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	//TODO: How does input come from front end??
	input := r.URL.Query().Get("query")
	if input == "" || len(input) > 1000 {
		errs.WriteError(w, fmt.Errorf("%w: invalid search length", errs.ErrInvalidUserInput))
		return
	}

	data, err := h.service.Search(ctx, input, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}
	data.Username = username

	data.Page = pagination.Page

	//TODO: what is Path?

	h.renderer.RenderPage(w, "searchresults.html", data)

	// in the repo layer the query will be.... what? need to figure that out first.

}
