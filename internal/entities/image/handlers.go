package image

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/internal/errs"
	"github.com/go-playground/validator/v10"
)

type ImageHandler struct {
	service   *ImageService
	validator *validator.Validate
}

const MaxUploadSize int64 = 10_000_000 // 10MB

func NewHandler(service *ImageService, validator *validator.Validate) *ImageHandler {
	return &ImageHandler{service: service, validator: validator}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: make a req struct instead of individual params
	threadID, err := strconv.Atoi(r.PostFormValue("thread_id"))
	if err != nil {
		errs.WriteError(w, err)
	}

	commentID, err := strconv.Atoi(r.PostFormValue("comment_id"))
	if err != nil {
		errs.WriteError(w, err)
	}

	// Limit upload file size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		errs.WriteError(w, err)
		return
	}

	imagePath, err := h.service.UploadImage(ctx, r, threadID, commentID)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"image_path": imagePath})
}
