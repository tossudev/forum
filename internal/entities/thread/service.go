package thread

import (
	
)

type ThreadService struct {
	repo *respository.Repo
}

func NewThreadService(r *respository.Repo) *ThreadService {
	return &ThreadService{repo: r}
}

func (s *ThreadService) GetThreadsByCategory() {
	
}