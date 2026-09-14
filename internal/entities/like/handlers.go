package like

import (
	"encoding/json"
	"net/http"

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
func (h *LikeHandler) LikeThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ThreadLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.WriteError(w, err)
		return
	}

	if err := h.service.LikeThread(ctx, req); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LikeHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CommentLikeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errs.WriteError(w, err)
		return
	}

	if err := h.service.LikeComment(ctx, req); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
