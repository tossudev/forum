package thread

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ThreadHandler struct {
	service   *ThreadService
	validator *validator.Validate
}

func NewHandler(service *ThreadService, validator *validator.Validate) *ThreadHandler {
	return &ThreadHandler{service: service, validator: validator}
}

func (h *ThreadHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {

}
