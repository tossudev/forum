package search

import (
	"errors"
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
	validator *validator.Validate
	renderer  *render.Renderer
}

func NewHandler(service *SearchService, validator *validator.Validate, renderer *render.Renderer) *SearchHandler {
	return &SearchHandler{service: service,
		renderer:  renderer,
		validator: validator}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	query := r.URL.Query()

	input := r.URL.Query().Get("search")

	if input == "" || len(input) > 1000 {
		errs.WriteError(w, fmt.Errorf("%w: invalid search length", errs.ErrInvalidUserInput))
		return
	}

	paginationSearch, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	//perform search, check for errs
	data, err := h.service.Search(ctx, input, paginationSearch)

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}
	data.Username = username

	//if err is 'no search results', send error to frontend and render search page
	//otherwise, WriteError as usual
	if err != nil {
		if errors.Is(err, errs.ErrSearchNoResults) {
			data.ResponseMsg = "No results found. Please try a different search term."
			w.WriteHeader(http.StatusBadRequest)
			h.renderer.RenderPage(w, "search.html", data)
			return
		} else {
			//errs.WriteError(w, err)
			w.WriteHeader(http.StatusInternalServerError)
			h.renderer.RenderPage(w, "error-500.html", nil)
			return
		}
	}

	//if no errors in search, then render search page and display search results

	data.Page = paginationSearch.Page

	data.ThreadPath = "Thread > Search"
	data.CommentPath = "Comment > Search"

	h.renderer.RenderPage(w, "search.html", data)
}
