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
	tab := r.URL.Query().Get("tab")

	var searchPageData SearchResult

	//get username if logged in
	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}
	searchPageData.Username = username

	//get pagination data from query
	paginationSearch, err := pagination.Parse(query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", searchPageData)
		return
	}

	//check for empty or too long search input
	searchPageData.InputQuery = input

	if searchPageData.InputQuery == "" || len(searchPageData.InputQuery) > 1000 {
		searchPageData.InputQuery = ""
		searchPageData.ResponseMsg = "Search query must be 1-1000 characters long."
		h.renderer.RenderPage(w, "search.html", searchPageData)
		return
	}

	//perform search, check for errs
	threads, threadsCount, comments, commentsCount, err := h.service.Search(ctx, input, paginationSearch)
	if tab == "threads" {
		comments = nil
	} else if tab == "comments" {
		threads = nil
	}

	searchPageData.Threads = threads
	searchPageData.Comments = comments

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

	searchPageData.ThreadTotalPages = pagination.CalcTotalPages(threadsCount, paginationSearch.PageSize)
	searchPageData.ThreadPageList = pagination.GetDisplayPageNums(paginationSearch.Page, searchPageData.ThreadTotalPages, 5) //5 is # pages displayed

	searchPageData.CommentTotalPages = pagination.CalcTotalPages(commentsCount, paginationSearch.PageSize)
	searchPageData.CommentPageList = pagination.GetDisplayPageNums(paginationSearch.Page, searchPageData.CommentTotalPages, 5)

	searchPageData.ThreadPage = paginationSearch.Page
	searchPageData.CommentPage = paginationSearch.Page

	searchPageData.ThreadPrevPage = searchPageData.ThreadPage - 1
	searchPageData.ThreadNextPage = searchPageData.ThreadPage + 1
	searchPageData.CommentPrevPage = searchPageData.CommentPage - 1
	searchPageData.CommentNextPage = searchPageData.CommentPage + 1

	searchPageData.ThreadPath = "Search > Thread"
	searchPageData.CommentPath = "Search > Comment"

	h.renderer.RenderPage(w, "search.html", searchPageData)
}
