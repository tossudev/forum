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
	"forum/internal/pagination"
	"forum/internal/render"
	"forum/internal/session"
	//"forum/internal/entities/user"
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
	query := r.URL.Query()

	paginationData, err := pagination.Parse(query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", ThreadsPage{})
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid category id", errs.ErrInvalidUserInput))
		return
	}

	threads, totalCount, err := h.service.GetByCategory(ctx, id, paginationData)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUserInput) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", ThreadsPage{})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		h.renderer.RenderPage(w, "error-500.html", ThreadsPage{})
		return
	}

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) // The visible pages in navigation

	data := ThreadsPage{
		// TODO: use category name instead of ID
		CategoryID: id,
		Path:       fmt.Sprintf("⌂ Home / CategoryID %d", id),
		Threads:    threads,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
		Username:   username,
	}

	h.renderer.RenderPage(w, "threads.html", data)
}

// TODO: if not being used by front end, make this func repo only??
func (h *ThreadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	//TODO: change _ to thread and send to front end
	_, err = h.service.GetByID(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ThreadHandler) CreateThreadPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := user.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
	}

	categoryID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || categoryID < 1 {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", CreateThreadsPage{})
		return
	}

	categoryName, err := h.service.GetCategoryName(ctx, categoryID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", CreateThreadsPage{})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		h.renderer.RenderPage(w, "error-500.html", CreateThreadsPage{})
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

	pageData := CreateThreadsPage{}
	user := user.GetUser(r)
	if user != nil {
		pageData.Username = user.Username
	}

	// Parse form before handling user not logged in error so that form data can be saved
	// and reinserted into the page templates when displaying errors.
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "error-400.html", pageData)
		slog.Error("thread creation fail", "err", err)
		return
	}

	categoryID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || categoryID < 1 {
		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "error-400.html", pageData)
		slog.Error("thread creation fail", "err", err)
		return
	}

	categoryName, err := h.service.GetCategoryName(ctx, categoryID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", CreateThreadsPage{})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		h.renderer.RenderPage(w, "error-500.html", CreateThreadsPage{})
	}

	thread := Thread{
		Title:      strings.TrimSpace(r.PostFormValue("title")),
		Body:       strings.TrimSpace(r.PostFormValue("body")),
		CategoryID: categoryID,
	}

	pageData.CategoryID = thread.CategoryID
	pageData.CategoryName = categoryName
	pageData.InputTitle = thread.Title
	pageData.InputBody = thread.Body

	// User not logged in
	if user == nil {
		pageData.OtherErr = "Please log in to create a thread." // Shows error without redirect so that the user doesn't lose their submission data
		w.WriteHeader(http.StatusUnauthorized)
		h.renderer.RenderPage(w, "thread-create.html", pageData)
		return
	}

	thread.AuthorID = user.ID
	pageData.Username = user.Username

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

	session := session.GetSession(r)
	if session == nil {
		errs.WriteError(w, fmt.Errorf("%w: not authenticated", errs.ErrUnauthorized))
		return
	}
	userID := session.UserID()

	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	err = h.service.Delete(ctx, threadID, userID)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ThreadHandler) DisplayThreadPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	data, err := h.service.DisplayThreadPage(ctx, id, username)
	if err != nil {
		errs.WriteError(w, err)
	}

	data.Username = username
	data.Path = "Thread > " + data.Thread.Title

	h.renderer.RenderPage(w, "thread.html", data)
}
