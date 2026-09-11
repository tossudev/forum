package user

type UserHandler struct {
	service *UserService
}

func NewHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}
