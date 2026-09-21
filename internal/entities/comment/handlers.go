package comment

import (
	"encoding/json"
	"fmt"
	"forum/internal/errs"
	"net/http"
	"strconv"

	"forum/internal/pagination"
	"forum/internal/session"

	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	service   *CommentService
	validator *validator.Validate
}

func NewHandler(service *CommentService, validator *validator.Validate) *CommentHandler {
	return &CommentHandler{service: service, validator: validator}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	if sess == nil {
		errs.WriteError(w, fmt.Errorf("%w: login required", errs.ErrUnauthorized))
		return
	}

	var input struct {
		Body string `json:"body" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid request body", errs.ErrInvalidUserInput))
		return
	}
	if err := h.validator.Struct(input); err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid input", errs.ErrInvalidUserInput))
		return
	}
	threadID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}
	req := Comment{
		Body:     input.Body,
		ThreadID: threadID,
		AuthorID: sess.UserID(),
	}
	_, err = h.service.Create(r.Context(), &req)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid comment id", errs.ErrInvalidUserInput))
		return
	}

	_, err = h.service.GetByID(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) GetByThread(w http.ResponseWriter, r *http.Request) {
	pagination, err := pagination.Parse(r.URL.Query())
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	_, err = h.service.GetByThread(r.Context(), id, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid id", errs.ErrInvalidUserInput))
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
