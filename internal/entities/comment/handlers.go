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
	sess, ok := session.GetSession(r.Context())
	if !ok || sess == nil { 
		errs.WriteError(w, fmt.Errorf("%w: login required", errs.ErrsUnauthorized))
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
		Body: input.Body,
		ThreadID: threadID,
		AuthorID: sess.UserID(),
	}
	_, err = h.service.Create(r.Context(), &req) //leftmost value left out, since the commented out section of json encoding (line 63)
	if err != nil {
		errs.WriteError(w, err)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(newComment) 
}

func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid comment id", errs.ErrInvalidUserInput))
		return
	}

	_, err = h.service.GetByID(ctx, id) //same thing for the comment, commented out the json encoding 
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // TODO: decision on what to return: statusOK or noContent?
	//json.NewEncoder(w).Encode(comment)
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

	//w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(comments)
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
