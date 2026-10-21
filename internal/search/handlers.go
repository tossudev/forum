package search

import (
	"errors"
	"net/http"

	"forum/internal/entities/user"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"
)

type SearchHandler struct {
	service  *SearchService
	renderer *render.Renderer
}

func NewHandler(service *SearchService, renderer *render.Renderer) *SearchHandler {
	return &SearchHandler{service: service,
		renderer: renderer}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	query := r.URL.Query()

	input := r.URL.Query().Get("search")
	tab := r.URL.Query().Get("tab")

	//get username if logged in
	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	//get pagination data from query
	paginationSearch, err := pagination.Parse(query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", SearchResult{Username: username})
		return
	}

	//check for empty or too long search input
	if input == "" || len(input) > 1000 {
		msg := "Search query must be 1-1000 characters long."
		h.renderer.RenderPage(w, "search.html", SearchResult{Username: username, ResponseMsg: msg})
		return
	}

	//perform search, check for errs
	threads, threadsCount, comments, commentsCount, err := h.service.Search(ctx, input, paginationSearch)
	if tab == "threads" {
		comments = nil
	} else if tab == "comments" {
		threads = nil
	}

	//if err is 'no search results', send error to frontend and render search page
	//otherwise, render error 500 page
	if err != nil {
		if errors.Is(err, errs.ErrSearchNoResults) {
			msg := "No results found. Please try a different search term."
			h.renderer.RenderPage(w, "search.html", SearchResult{Username: username, ResponseMsg: msg})
			return
		} else if errors.Is(err, errs.ErrInvalidUserInput) || errors.Is(err, errs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", SearchResult{Username: username})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			h.renderer.RenderPage(w, "error-500.html", SearchResult{Username: username})
			return
		}
	}

	//if no errors in search, then render search page and display search results

	//pagination for threads
	threadTotalPages := pagination.CalcTotalPages(threadsCount, paginationSearch.PageSize)
	threadPageList := pagination.GetDisplayPageNums(paginationSearch.Page, threadTotalPages, 5) //5 is #pages displayed

	//pagination for comments
	commentTotalPages := pagination.CalcTotalPages(commentsCount, paginationSearch.PageSize)
	commentPageList := pagination.GetDisplayPageNums(paginationSearch.Page, commentTotalPages, 5)

	threadPageNav := pagination.PaginationNav{
		Path:       "Search > Thread",
		Page:       paginationSearch.Page,
		TotalPages: threadTotalPages,
		PageList:   threadPageList,
		PrevPage:   paginationSearch.Page - 1,
		NextPage:   paginationSearch.Page + 1,
	}

	commentPageNav := pagination.PaginationNav{
		Path:       "Search > Comment",
		Page:       paginationSearch.Page,
		TotalPages: commentTotalPages,
		PageList:   commentPageList,
		PrevPage:   paginationSearch.Page - 1,
		NextPage:   paginationSearch.Page + 1,
	}

	pageData := SearchResult{
		InputQuery:     input,
		Threads:        threads,
		Comments:       comments,
		ThreadPageNav:  threadPageNav,
		CommentPageNav: commentPageNav,
		Username:       username,
	}

	h.renderer.RenderPage(w, "search.html", pageData)
}
