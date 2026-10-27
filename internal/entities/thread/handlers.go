package thread

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/entities/user"
	"forum/internal/errs"
	"forum/internal/models"
	"forum/internal/pagination"
	"forum/internal/render"
	"forum/internal/validate"
)

type ThreadHandler struct {
	service  *ThreadService
	renderer *render.Renderer
}

func NewHandler(service *ThreadService, renderer *render.Renderer) *ThreadHandler {
	return &ThreadHandler{
		service:  service,
		renderer: renderer,
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
		Title:      validate.TrimEmpty(r.PostFormValue("title")),
		Body:       validate.TrimEmpty(r.PostFormValue("body")),
		CategoryID: categoryID,
	}

	pageData := CreateThreadsPage{
		CategoryID:   thread.CategoryID,
		CategoryName: categoryName,
		InputTitle:   thread.Title,
		InputBody:    thread.Body,
		Username:     username,
	}

	fieldErrors := make(map[string]string)

	// User not logged in: Shows error without redirect so that the user doesn't lose their submission data
	if user == nil {
		fieldErrors["OtherErr"] = "Please log in to create a thread."
		pageData.FieldErrors = fieldErrors
		w.WriteHeader(http.StatusUnauthorized)
		h.renderer.RenderPage(w, "thread-create.html", pageData)
		return
	}

	thread.AuthorID = user.ID

	threadID, err := h.service.Create(ctx, &thread)
	if err != nil {
		fieldErrors, statusCode := getThreadCreationFormFieldErrors(err)
		pageData.FieldErrors = fieldErrors
		w.WriteHeader(statusCode)
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

	user := user.GetUser(r)
	username := ""
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

	params := ThreadPageParams{
		ThreadID:       threadID,
		PaginationData: paginationData,
		Username:       username,
		Path:           r.URL.Path,
	}
	pageData, totalCount, err := h.service.CollectPageData(ctx, params)
	if err != nil {
		if errors.Is(err, pagination.ErrPageOverTotal) {
			totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
			http.Redirect(w, r, r.URL.Path+"?page="+strconv.Itoa(totalPages), http.StatusSeeOther) // Redirect to last page - good UX when a user deletes the only comment on the last page (otherwise they see 404)
			return
		}
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: username})
		return
	}
	h.renderer.RenderPage(w, "thread.html", pageData)
}
