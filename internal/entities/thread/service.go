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

	threads, err := s.repo.GetByCategory(ctx, id, pagination)

	for i, thread := range threads {
		var errFormat error
		threads[i].DateCreated, errFormat = config.FormatTime(thread.DateCreated)
		if errFormat != nil {
			return nil, errFormat
		}
	}

	return threads, err
}

func (s *ThreadService) GetByID(ctx context.Context, id int) (*Thread, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ThreadService) Create(ctx context.Context, thread *Thread) error {
	now := time.Now().Format(config.TimeFormat)

	thread.DateCreated = now

	return s.repo.Create(ctx, thread)
}

func (s *ThreadService) Delete(ctx context.Context, threadID int, userID int) error {
	if threadID <= 0 {
		return fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}
	return s.repo.Delete(ctx, threadID, userID)
}

func (s *ThreadService) DisplayThreadPage(ctx context.Context, threadID int, username string) (*ThreadData, error) {
	if threadID <= 0 {
		return nil, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}

	data, err := s.repo.DisplayThreadPage(ctx, threadID, username)

	var errFormat error
	data.Thread.DateCreated, errFormat = config.FormatTime(data.Thread.DateCreated)
	if errFormat != nil {
		return nil, errFormat
	}

	for i, comment := range data.Comments {
		data.Comments[i].DateCreated, errFormat = config.FormatTime(comment.DateCreated)
		if errFormat != nil {
			return nil, errFormat
		}
	}

	return data, err
}
