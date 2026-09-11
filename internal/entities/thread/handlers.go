package thread

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		//TODO: handle this error
		return
	}

	threads, err := h.service.GetByCategory(id)
	if err != nil {
		//TODO: handle this error
		return
	}

	//TODO: find out what front end needs this result to do/look like
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(threads)
}

func (h *ThreadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		//TODO: handle this error
		return
	}

	thread, err := h.service.GetByID(id)
	if err != nil {
		//TODO: handle this error
		fmt.Println(err)
	}

	//TODO: find out what front end needs this result to do/look like
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(thread)
}

func (h *ThreadHandler) Create(w http.ResponseWriter, r *http.Request) {
	//TODO: agree/confirm form values
	//TODO: use these when form exists
	/*category := r.FormValue("category")
	title := r.FormValue("title")
	body := r.FormValue("body")*/

	//FOR TESTING REMOVE THIS WHEN FRONT END EXISTS...
	var threadRequest *Thread
	err := json.NewDecoder(r.Body).Decode(&threadRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	//...UNTIL HERE

	thread, err := h.service.Create(threadRequest)
	if err != nil {
		//TODO: handle this error
		fmt.Println(err)
		return
	}

	//TODO: find out what front end needs this result to do/look like
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(thread)

}
