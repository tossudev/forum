package comment

import (
	//	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	service   *CommentService
	validator *validator.Validate
}

func NewHandler(service *CommentService, validator *validator.Validate) *CommentHandler {
	return &CommentHandler{service: service, validator: validator}
}

func Create(w http.ResponseWriter, r *http.Request) {

	return
}

func GetByThread(w http.ResponseWriter, r *http.Request) {
	return
}
