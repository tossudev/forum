package thread

import (

)

type ThreadHandler struct {
	service   *ThreadService
	validator *validator.Validate
}

func NewHandler(service *ThreadService, validator *validator.Validate) {
	return &ThreadHandler{service: service, validator: validator}
}

func (h *ThreadHandler) GetByCategory(w http.ResponseWriter, r *http.Response) {

}