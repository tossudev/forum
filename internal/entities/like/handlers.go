package like

import (
	"encoding/json"
	"net/http"

	"forum/internal/entities/user"
	"forum/internal/errs"

	"github.com/go-playground/validator/v10"
)

type LikeHandler struct {
	service  *LikeService
	validate *validator.Validate
}

func NewHandler(service *LikeService, validate *validator.Validate) *LikeHandler {
	return &LikeHandler{service: service, validate: validate}
}

// TODO: handle errors
func (h *LikeHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// type Request struct {
	// 	CommentID int  `json:"id"`
	// 	Like      bool `json:"like"`
	// }
	// var req Request
	var req CommentLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//TODO: show 400 page
		errs.WriteError(w, err)
		return
	}

	if user := user.GetUser(r); user != nil {
		req.UserID = user.ID
	} else {
		//TODO: Show user-friendly error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.service.LikeComment(ctx, req); err != nil {
		//TODO: Show user-friendly error
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO: handle errors
func (h *LikeHandler) LikeThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ThreadLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//TODO: show 400 page
		errs.WriteError(w, err)
		return
	}

	if user := user.GetUser(r); user != nil {
		req.UserID = user.ID
	} else {
		//TODO: Show user-friendly error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.service.LikeThread(ctx, req); err != nil {
		//TODO: Show user-friendly error
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
