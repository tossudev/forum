package like

import (
	"encoding/json"
	"net/http"

	"forum/internal/entities/user"
	"forum/internal/errs"
)

type LikeHandler struct {
	service *LikeService
}

func NewHandler(service *LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

func (h *LikeHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CommentLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	if user := user.GetUser(r); user != nil {
		req.UserID = user.ID
	} else {
		errs.WriteError(w, errs.ErrUnauthorized)
		return
	}

	if err := h.service.LikeComment(ctx, req); err != nil {
		errs.WriteError(w, err)
		return
	}
}

func (h *LikeHandler) LikeThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ThreadLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.WriteError(w, errs.ErrInvalidUserInput)
		return
	}

	if user := user.GetUser(r); user != nil {
		req.UserID = user.ID
	} else {
		errs.WriteError(w, errs.ErrUnauthorized)
		return
	}

	if err := h.service.LikeThread(ctx, req); err != nil {
		errs.WriteError(w, err)
		return
	}
}
