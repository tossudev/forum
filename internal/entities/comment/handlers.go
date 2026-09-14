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
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&test)
	if err != nil {
		//handle this error
	}

	comment := Comment{
		Body: test.Body,
	}

	err = h.service.Create(r.Context(), &comment)
	if err != nil {
		//handle this error
	}

	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(test)

	return
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
