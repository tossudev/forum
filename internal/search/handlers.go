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

	searchQuery := r.URL.Query().Get("search")
	tab := r.URL.Query().Get("tab")

	//get username if logged in
	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	//get pagination data from query
	paginationData, err := pagination.Parse(query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", SearchResultPage{Username: username})
		return
	}

	//check for empty or too long search searchQuery
	if searchQuery == "" || len(searchQuery) > 1000 {
		msg := "Search query must be 1-1000 characters long."
		h.renderer.RenderPage(w, "search.html", SearchResultPage{Username: username, ResponseMsg: msg})
		return
	}

	//perform search, check for errs
	searchResult, err := h.service.Search(ctx, searchQuery, tab, paginationData)
	paginationData.Page = searchResult.Page

	//if err is 'no search results', send error to frontend and render search page
	//otherwise, render error 500 page
	if err != nil {
		if errors.Is(err, errs.ErrSearchNoResults) {
			msg := "No results found. Please try a different search term."
			h.renderer.RenderPage(w, "search.html", SearchResultPage{Username: username, ResponseMsg: msg})
			return
		} else if errors.Is(err, errs.ErrInvalidUserInput) || errors.Is(err, errs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", SearchResultPage{Username: username})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			h.renderer.RenderPage(w, "error-500.html", SearchResultPage{Username: username})
			return
		}
	}

	//if no errors in search, then render search page and display search results

	//pagination for threads
	threadTotalPages := pagination.CalcTotalPages(searchResult.ThreadsCount, paginationData.PageSize)
	threadPageList := pagination.GetDisplayPageNums(paginationData.Page, threadTotalPages, 5) //5 is #pages displayed

	//pagination for comments
	commentTotalPages := pagination.CalcTotalPages(searchResult.CommentsCount, paginationData.PageSize)
	commentPageList := pagination.GetDisplayPageNums(paginationData.Page, commentTotalPages, 5)

	threadPageNav := pagination.PaginationNav{
		Page:       paginationData.Page,
		TotalPages: threadTotalPages,
		PageList:   threadPageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}

	commentPageNav := pagination.PaginationNav{
		Page:       paginationData.Page,
		TotalPages: commentTotalPages,
		PageList:   commentPageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}

	pageData := SearchResultPage{
		ThreadPath:     "Search > Thread",
		CommentPath:    "Search > Comment",
		InputQuery:     searchQuery,
		Threads:        searchResult.Threads,
		Comments:       searchResult.Comments,
		ThreadPageNav:  threadPageNav,
		CommentPageNav: commentPageNav,
		Username:       username,
	}

	h.renderer.RenderPage(w, "search.html", pageData)
}
