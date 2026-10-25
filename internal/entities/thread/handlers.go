package thread

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"forum/internal/entities/user"
	"forum/internal/errs"
	"forum/internal/models"
	"forum/internal/pagination"
	"forum/internal/render"
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

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	threads, totalCount, categoryName, err := h.service.GetByCategory(ctx, id, paginationData)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: username})
		return
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) // The visible pages in navigation

	pageNav := pagination.PaginationNav{
		Path:       r.URL.Path,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}

	data := ThreadsPage{
		CategoryID: id,
		Breadcrumbs: []models.Breadcrumb{
			models.Breadcrumb{
				Name: "Home",
				URL:  "/",
			},
			models.Breadcrumb{
				Name: categoryName,
				URL:  "/categories/" + strconv.Itoa(id),
			},
		},
		Threads:  threads,
		PageNav:  pageNav,
		Username: username,
	}

	h.renderer.RenderPage(w, "threads.html", data)
}

func (h *ThreadHandler) CreateThreadPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var username string
	user := user.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}
	username = user.Username

	categoryID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || categoryID < 1 {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	categoryName, err := h.service.GetCategoryName(ctx, categoryID)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: username})
		return
	}

	pageData := CreateThreadsPage{
		CategoryID:   categoryID,
		CategoryName: categoryName,
		Username:     user.Username,
	}

	h.renderer.RenderPage(w, "thread-create.html", pageData)
}

func (h *ThreadHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	// Parse form before handling user not logged in error so that form data can be saved
	// and reinserted into the page templates when displaying errors.
	if err := r.ParseForm(); err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	categoryID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || categoryID < 1 {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	categoryName, err := h.service.GetCategoryName(ctx, categoryID)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: username})
	}

	thread := Thread{
		Title:      strings.TrimSpace(r.PostFormValue("title")),
		Body:       strings.TrimSpace(r.PostFormValue("body")),
		CategoryID: categoryID,
	}

	pageData := CreateThreadsPage{
		CategoryID:   thread.CategoryID,
		CategoryName: categoryName,
		InputTitle:   thread.Title,
		InputBody:    thread.Body,
		Username:     username,
	}

	// User not logged in: Shows error without redirect so that the user doesn't lose their submission data
	if user == nil {
		pageData.OtherErr = "Please log in to create a thread."
		w.WriteHeader(http.StatusUnauthorized)
		h.renderer.RenderPage(w, "thread-create.html", pageData)
		return
	}

	thread.AuthorID = user.ID

	threadID, err := h.service.Create(ctx, &thread)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidThreadTitle) {
			pageData.TitleErr = "Thread title must be 1-100 characters."
		}

		if errors.Is(err, errs.ErrInvalidThreadBody) {
			pageData.BodyErr = "Thread body must be 1-5000 characters."
		}

		if !errors.Is(err, errs.ErrInvalidThreadTitle) && !errors.Is(err, errs.ErrInvalidThreadBody) {
			pageData.OtherErr = "Unable to create thread. Please try again later."
			w.WriteHeader(http.StatusInternalServerError)
			h.renderer.RenderPage(w, "thread-create.html", pageData)
			slog.Error("thread creation fail", "err", err)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "thread-create.html", pageData)
		return
	}

	// TODO: On successful submission, redirect to newly created thread
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusSeeOther)
}

func (h *ThreadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := user.GetUser(r)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	err = h.service.Delete(ctx, threadID, user.ID)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ThreadHandler) DisplayThreadPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	if err := paginationData.Validate(); err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || threadID < 1 {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: username})
		return
	}

	threadDisplay, commentsDisplay, totalCount, err := h.service.DisplayThreadPage(ctx, threadID, username, paginationData)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: username})
		return
	}

	catID := threadDisplay.Thread.CategoryID
	categoryName, err := h.service.GetCategoryName(ctx, catID)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: username})
		return
	}

	breadcrumbs := []models.Breadcrumb{
		models.Breadcrumb{
			Name: "Home",
			URL:  "/",
		},
		models.Breadcrumb{
			Name: categoryName,
			URL:  "/categories/" + strconv.Itoa(catID),
		},
		models.Breadcrumb{
			Name: threadDisplay.Thread.Title,
			URL:  "/threads/" + strconv.Itoa(threadID),
		},
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) // The visible pages in navigation
	pageNav := pagination.PaginationNav{
		Path:       r.URL.Path,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}

	pageData := ThreadCommentsPage{
		Breadcrumbs:   breadcrumbs,
		PageNav:       pageNav,
		Username:      username,
		ThreadDisplay: *threadDisplay,
		Comments:      commentsDisplay,
	}

	h.renderer.RenderPage(w, "thread.html", pageData)
}
