package thread

import (

)

type ThreadHandler struct {
	service   *ThreadService
	validator *validator.Validate
}

func NewThreadHandler(service *ThreadService, validator *validator.Validate) {
	return &ThreadHandler{service: service, validator: validator}
}

func (h *ThreadHandler) GetThreadsByCategory(w http.ResponseWriter, r *http.Response) {

}