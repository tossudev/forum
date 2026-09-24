package image

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"forum/internal/errs"
	"github.com/go-playground/validator/v10"
)

type ImageHandler struct {
	service   *ImageService
	validator *validator.Validate
}

const MaxUploadSize int64 = 3_000_000 // 3MB

func NewHandler(service *ImageService, validator *validator.Validate) *ImageHandler {
	return &ImageHandler{service: service, validator: validator}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Limit upload file size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			errs.WriteError(w, errs.ErrFileTooLarge)
			return
		}

		errs.WriteError(w, err)
		return
	}

	req := ImageRequest{}
	var err error

	req.ThreadID, err = strconv.Atoi(r.PostFormValue("thread_id"))
	if err != nil {
		errs.WriteError(w, err)
	}

	req.CommentID, err = strconv.Atoi(r.PostFormValue("comment_id"))
	if err != nil {
		errs.WriteError(w, err)
	}

	imagePath, err := h.service.UploadImage(ctx, r, req)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"image_path": imagePath})
}
