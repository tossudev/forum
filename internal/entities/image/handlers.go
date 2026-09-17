package image

import (
	"net/http"

	"forum/internal/errs"
	"github.com/go-playground/validator/v10"
)

type ImageHandler struct {
	service   *ImageService
	validator *validator.Validate
}

const MaxUploadSize int64 = 1024 * 1024 // 1MB

func NewHandler(service *ImageService, validator *validator.Validate) *ImageHandler {
	return &ImageHandler{service: service, validator: validator}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Limit upload file size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		errs.WriteError(w, err)
		return
	}

	if err := h.service.UploadImage(ctx, r); err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
