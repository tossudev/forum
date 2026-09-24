package user

import (
	"encoding/json"
	"fmt"
	"forum/internal/errs"
	"forum/internal/password"
	"forum/internal/render"
	"forum/internal/session"
	"net/http"
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

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid request body", errs.ErrInvalidUserInput))
		return
	}

	user := User{
		Username: input.Username,
		Email:    input.Email,
		Password: password.Password{},
		RoleID:   1,
	}

	if err := h.service.RegisterUser(ctx, &user, input.Password); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Welcome to Literary Lions, %s!", user.Username)})
}

func (h *UserHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type PageData struct {
		User string
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
