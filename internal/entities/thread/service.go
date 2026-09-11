package thread

import (
)

type ThreadService struct {
	repo *ThreadRepository
}

func NewService(r *ThreadRepository) *ThreadService {
	return &ThreadService{repo: r}
}

func (s *ThreadService) GetByCategory() {
	
}
