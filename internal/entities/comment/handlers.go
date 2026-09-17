package comment

import (
	"encoding/json"
	"fmt"
	"forum/internal/errs"
	"net/http"
	"strconv"

	"forum/internal/pagination"

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


	threadIDstring := r.PathValue("id")
	threadID, err := strconv.Atoi(threadIDstring)	
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput))
		return
	}

	req := Comment{
		Body: input.Body,
		ThreadID: threadID,
		AuthorID: 1, 		//TODO: we get the author id from sessions. Will be implemented later.
	}

	newComment, err := h.service.Create(r.Context(), &req)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}

func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	//endpoint: mux.HandleFunc(GET comments/{id})

	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		errs.WriteError(w, fmt.Errorf("%w: invalid comment id", errs.ErrInvalidUserInput))
		return
	}

	comment, err := h.service.GetByID(ctx, id)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetByThread(w http.ResponseWriter, r *http.Request) {

	//endpoint: mux.HandleFunc(GET threads/{id}/comments)

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

	comments, err := h.service.GetByThread(r.Context(), id, pagination)
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}


func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO
	//endpoint: mux.HandleFunc(DELETE comments/{id})

	return 
}





