package thread

import (
	
)

type ThreadService struct {
	repo *respository.Repo
}

func NewService(r *respository.Repo) *ThreadService {
	return &ThreadService{repo: r}
}

func (s *ThreadService) GetByCategory() {
	
}