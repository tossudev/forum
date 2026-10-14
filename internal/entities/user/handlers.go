package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"forum/internal/errs"
	"forum/internal/password"
	"forum/internal/render"
	"forum/internal/session"
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

	user := &User{
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: password.Password{},
		RoleID:   1,
	}

	user, err := h.service.RegisterUser(ctx, user, r.PostFormValue("password"))
	if err != nil {
		pageData := RegistrationPage{Username: r.PostFormValue("username"), Email: r.PostFormValue("email")} // service.RegisterUser returns a nil user on error, so values are taken from the form

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
			pageData.PasswordErr = "Passwords must be at least 8 characters long."
		}

		if !errors.Is(err, errs.ErrInvalidUserInput) && !errors.Is(err, errs.ErrDuplicateUsername) && !errors.Is(err, errs.ErrDuplicateEmail) {
			pageData.ServerErr = "Sorry, we're unable to process user registration at the moment. Please try again later."
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

	type PageData struct {
		Username string
	}

	h.renderer.RenderPage(w, "login.html", PageData{""})
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

	input := CredentialsSubmission{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	userID, err := h.service.Authenticate(ctx, input)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	// Create new session
	if err := h.sm.Login(userID, w, r); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged in successfully"})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sm.Logout(w, r); err != nil {
		errs.WriteError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
