package user

import (
	"encoding/json"
	"fmt"
	"forum/internal/errs"
	"forum/internal/password"
	"net/http"
)

type UserHandler struct {
	service *UserService
}

func NewHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
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
	}

	if err := h.service.RegisterUser(ctx, &user, input.Password); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprint("Welcome to Literary Lions, %s!", user.Username)})
}
