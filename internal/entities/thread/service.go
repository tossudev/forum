package thread

import (
	
)

type ThreadService struct {
	repo *respository.Repo
	validate *validator.Validate
}

func NewThreadService(r *respository.Repo, v *validator.Validate) *ThreadService {
	return &ThreadService{repo: r, validate: v}
}
