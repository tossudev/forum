package thread

import (
	"context"
	"fmt"
	"time"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
)

type ThreadService struct {
	repo *ThreadRepository
}

func NewService(r *ThreadRepository) *ThreadService {
	return &ThreadService{repo: r}
}

func (s *ThreadService) GetByCategory(ctx context.Context, id int, pagination pagination.Pagination) ([]Thread, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid category id", errs.ErrInvalidUserInput)
	}

	return s.repo.GetByCategory(ctx, id, pagination)
}

func (s *ThreadService) GetByID(ctx context.Context, id int) (*Thread, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ThreadService) Create(ctx context.Context, thread *Thread) (*Thread, error) {

	now := time.Now().Format(config.TimeFormat)

	thread.DateCreated = now

	return s.repo.Create(ctx, thread)
}
