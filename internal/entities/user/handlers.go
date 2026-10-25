package user

import (
	"errors"
	"fmt"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/password"
	"forum/internal/render"
	"forum/internal/session"
	"log/slog"
	"net/http"
	"strings"
)

type UserHandler struct {
	service  *UserService
	sm       *session.SessionManager
	renderer *render.Renderer
}

func NewHandler(service *UserService, sm *session.SessionManager, renderer *render.Renderer) *UserHandler {
	return &UserHandler{
		service:  service,
		sm:       sm,
		renderer: renderer,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userSession := session.GetSession(r)
	if userSession == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	user, err := h.service.GetUserByID(ctx, userSession.UserID())
	if err != nil {
		errs.WriteError(w, fmt.Errorf("get user by id: %w", err))
	}

	// TODO: figure out roles
	data := ProfilePage{
		Username:    user.Username,
		Role:        "Member",
		DateCreated: user.CreatedAt.Format("02 Jan 2006"),
	}

	h.renderer.RenderPage(w, "profile.html", data)
}

func (h *UserHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.renderer.RenderPage(w, "register.html", RegistrationPage{})
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	inputUsername := strings.TrimSpace(r.PostFormValue("username"))
	inputEmail := strings.TrimSpace(r.PostFormValue("email"))

	user := &User{
		Username: inputUsername,
		Email:    inputEmail,
		Password: password.Password{},
		RoleID:   1,
	}

	user, err := h.service.RegisterUser(ctx, user, r.PostFormValue("password"))
	if err != nil {
		pageData := RegistrationPage{InputUsername: inputUsername, InputEmail: inputEmail} // service.RegisterUser returns a nil user on error, so values are taken from the form

		switch { // Username field errors
		case errors.Is(err, errs.ErrInvalidUsername):
			pageData.UsernameErr = "Username must be 2-20 characters, start with a letter, and only contain letters, numbers, or underscores (_)."
		case errors.Is(err, errs.ErrDuplicateUsername):
			pageData.UsernameErr = "Username already taken."
		}

		switch { // Email field errors
		case errors.Is(err, errs.ErrInvalidEmail):
			pageData.EmailErr = "Enter a valid email address, such as example@email.com."
		case errors.Is(err, errs.ErrDuplicateEmail):
			pageData.EmailErr = "Email already taken."
		}

		// Password field errors
		if errors.Is(err, errs.ErrInvalidPassword) {
			pageData.PasswordErr = "Password must be 8-128 characters."
		}

		if !errors.Is(err, errs.ErrInvalidUserInput) && !errors.Is(err, errs.ErrDuplicateUsername) && !errors.Is(err, errs.ErrDuplicateEmail) {
			pageData.OtherErr = "Sorry, we're unable to process user registration at the moment. Please try again later."
			slog.Error("internal server error:", "err", err)
		}

		w.WriteHeader(http.StatusBadRequest) //TODO: determine status code based on error
		h.renderer.RenderPage(w, "register.html", pageData)
		//errs.WriteError(w, err)
		return
	}

	// If registration successful, login user and redirect to home
	// Create new session
	if err := h.sm.Login(user.ID, w, r); err != nil {
		errs.WriteError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.renderer.RenderPage(w, "login.html", LoginPage{})
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	inputEmail := strings.TrimSpace(r.PostFormValue("email"))
	inputPassword := r.PostFormValue("password")

	input := LoginSubmission{
		Email:    inputEmail,
		Password: inputPassword,
	}

	userID, err := h.service.Authenticate(ctx, input)
	if err != nil {
		pageData := LoginPage{Email: inputEmail}

		switch {
		case errors.Is(err, errs.ErrInvalidUserInput):
			if errors.Is(err, errs.ErrInvalidEmail) {
				pageData.EmailErr = "Enter a valid email address, such as example@email.com."
			}
			if errors.Is(err, errs.ErrInvalidPassword) {
				pageData.PasswordErr = "Password must be 8-128 characters."
			}
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, errs.ErrUnauthorized):
			pageData.OtherErr = "Invalid credentials."
			w.WriteHeader(http.StatusUnauthorized)
		default:
			pageData.OtherErr = "Sorry, we're unable to process user login at the moment. Please try again later."
			slog.Error("internal server error:", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		h.renderer.RenderPage(w, "login.html", pageData)
		return
	}

	// Create new session
	if err := h.sm.Login(userID, w, r); err != nil {
		errs.WriteError(w, err)
		return
	}

	// Redirect to home on successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sm.Logout(w, r); err != nil {
		slog.Error("logout error", "err", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ProfileThreadsPage gets all threads posted by a user and renders the data on the user's profile page.
func (h *UserHandler) ProfileThreadsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	type ThreadsPage struct {
		Threads  []ProfileThread
		PageNav  pagination.PaginationNav
		Username string
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", ThreadsPage{Username: user.Username})
		return
	}

	if paginationData.Validate() != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "error-404.html", ThreadsPage{Username: user.Username})
		return
	}

	// Get all the threads posted by a user
	threads, totalCount, err := h.service.GetThreadsByUserID(ctx, user.ID, paginationData)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUserInput) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", ThreadsPage{Username: user.Username})
			return
		}
		slog.Error("internal server error:", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.renderer.RenderPage(w, "error-500.html", ThreadsPage{Username: user.Username})
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
	pageData := ThreadsPage{
		Threads:  threads,
		PageNav:  pageNav,
		Username: user.Username,
	}

	h.renderer.RenderPage(w, "profile-threads.html", pageData)
}

// ProfileLikedThreadsPage gets all liked threads by a user and renders the data on the user's profile page.
func (h *UserHandler) ProfileLikedThreadsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	type LikesPage struct {
		Threads  []ProfileThread
		PageNav  pagination.PaginationNav
		Username string
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.renderer.RenderPage(w, "error-404.html", LikesPage{Username: user.Username})
		return
	}

	if paginationData.Validate() != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.renderer.RenderPage(w, "error-404.html", LikesPage{Username: user.Username})
		return
	}

	// Get all the liked threads posted by a user
	threads, totalCount, err := h.service.GetLikedThreadsByUserID(ctx, user.ID, paginationData)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidUserInput) {
			w.WriteHeader(http.StatusNotFound)
			h.renderer.RenderPage(w, "error-404.html", LikesPage{Username: user.Username})
			return
		}
		slog.Error("internal server error:", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.renderer.RenderPage(w, "error-500.html", LikesPage{Username: user.Username})
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
	pageData := LikesPage{
		Threads:  threads,
		PageNav:  pageNav,
		Username: user.Username,
	}

	h.renderer.RenderPage(w, "profile-liked-threads.html", pageData)
}
