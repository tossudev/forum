package thread

import (
	"errors"
)

type ThreadService struct {
	repo *ThreadRepository
}

func NewService(r *ThreadRepository) *ThreadService {
	return &ThreadService{repo: r}
}

func (s *ThreadService) GetByCategory(id int) ([]Thread, error) {
	if id <= 0 {
		return nil, errors.New("invalidCategoryID")
	}

	threads, err := s.repo.GetByCategory(id)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (s *ThreadService) GetByID(id int) (*Thread, error) {
	if id <= 0 {
		return nil, errors.New("invalidThreadID")
	}

	thread, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (s *ThreadService) Create(threadRequest *Thread) (*Thread, error) {

	thread, err := s.repo.Create(threadRequest)
	if err != nil {
		return nil, err
	}

	return thread, nil
}
