package comment

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	var test struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&test)
	if err != nil {
		//handle this error
	}

	req := Comment{
		Body: test.Body,
	}

	newComment, err := h.service.Create(r.Context(), &req)
	if err != nil {
		//handle this error
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newComment)

	return
}

func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		//handle error
	}

	comment, err := h.service.GetByID(ctx, id)
	if err != nil {
		//handle error
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetByThread(w http.ResponseWriter, r *http.Request) {

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		//handle this error
	}

	comments, err := h.service.GetByThread(r.Context(), id)
	if err != nil {
		//handle this error
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)

	return
}
