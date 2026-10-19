package search

import (
	"errors"
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

	var searchPageData SearchResult

	//get username if logged in
	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}
	searchPageData.Username = username

	//check for empty or too long search input
	if input == "" || len(input) > 1000 {
		searchPageData.InputQuery = ""
		searchPageData.ResponseMsg = "Search query must be 1-1000 characters long."
		h.renderer.RenderPage(w, "search.html", searchPageData)
		return
	}

	//get pagination data from query
	paginationSearch, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	//perform search, check for errs
	threads, comments, err := h.service.Search(ctx, input, paginationSearch)

	searchPageData.Threads = threads
	searchPageData.Comments = comments
	searchPageData.InputQuery = input

	//if err is 'no search results', send error to frontend and render search page
	//otherwise, render error 500 page
	if err != nil {
		if errors.Is(err, errs.ErrSearchNoResults) {
			searchPageData.ResponseMsg = "No results found. Please try a different search term."
			h.renderer.RenderPage(w, "search.html", searchPageData)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			h.renderer.RenderPage(w, "error-500.html", nil)
			return
		}
	}

	//if no errors in search, then render search page and display search results

	searchPageData.Page = paginationSearch.Page

	searchPageData.ThreadPath = "Thread > Search"
	searchPageData.CommentPath = "Comment > Search"

	h.renderer.RenderPage(w, "search.html", searchPageData)
}
