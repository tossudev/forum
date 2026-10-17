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
	pagination, err := pagination.Parse(query)
	if err != nil {
		errs.WriteError(w, err)
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid category id", errs.ErrInvalidUserInput))
		return
	}

	threads, err := h.service.GetByCategory(ctx, id, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	var username string
	user := user.GetUser(r)
	if user != nil {
		username = user.Username
	}

	data := ThreadsPage{
		// TODO: use category name instead of ID
		CategoryID: id,
		Path:       fmt.Sprintf("⌂ Home / CategoryID %d", id),
		Threads:    threads,
		Page:       pagination.Page,
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

func (h *ThreadHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	user := user.GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
	}

	categoryID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		//TODO: display error message that the category doesn't exist
	}

	//TODO: Get category name
	pageData := CreateThreadsPage{
		CategoryID: categoryID,
		// CategoryName: ,
		Username: user.Username,
	}

	h.renderer.RenderPage(w, "thread-create.html", pageData)
}

func (h *ThreadHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session := session.GetSession(r)
	if session == nil {
		//TODO: show user friendly error
		errs.WriteError(w, fmt.Errorf("%w: not authenticated", errs.ErrUnauthorized))
		return
	}
	authorID := session.UserID()

	if err := r.ParseForm(); err != nil {
		//TODO: show user friendly error
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	categoryID, err := strconv.Atoi(r.PostFormValue("categoryID"))
	if err != nil {
		//TODO: show user friendly error
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	thread := Thread{
		Title:      strings.TrimSpace(r.PostFormValue("title")),
		Body:       strings.TrimSpace(r.PostFormValue("body")),
		AuthorID:   authorID,
		CategoryID: categoryID,
	}

	threadID, err := h.service.Create(ctx, &thread)
	if err != nil {
		pageData := CreateThreadsPage{CategoryID: thread.CategoryID, InputTitle: thread.Title, InputBody: thread.Body}

		if errors.Is(err, errs.ErrInvalidThreadTitle) {
			pageData.TitleErr = "Thread title must be 1-100 characters."
		}

		if errors.Is(err, errs.ErrInvalidThreadBody) {
			pageData.BodyErr = "Thread body must be 1-5000 characters."
		}

		if !errors.Is(err, errs.ErrInvalidThreadTitle) && errors.Is(err, errs.ErrInvalidThreadBody) {
			pageData.OtherErr = "Unable to create thread. Please try again later."
			slog.Error("thread creation fail", "err", err)
		}

		h.renderer.RenderPage(w, "thread-create.html", pageData)
		return
	}

	// TODO: On successful submission, redirect to newly created thread
	http.Redirect(w, r, "threads/"+strconv.Itoa(threadID), http.StatusSeeOther)
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
