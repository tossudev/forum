package thread

import (
	"context"
	"fmt"
	"time"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/render"
)

type ThreadService struct {
	repo *ThreadRepository
}

func NewService(r *ThreadRepository) *ThreadService {
	return &ThreadService{repo: r}
}

func (s *ThreadService) GetByCategory(ctx context.Context, id int, pagination pagination.Pagination, renderer *render.Renderer) ([]Thread, error) {
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

		threads[i].BodyHTML, err = renderer.ParseHTML(threads[i].Body)
		if err != nil {
			return nil, fmt.Errorf("GetByCategory ParseHTML err: %w", err)
		}

		threads[i].TitleHTML, err = renderer.ParseHTML(threads[i].Title)
		if err != nil {
			return nil, fmt.Errorf("GetByCategory ParseHTML err: %w", err)
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

func (s *ThreadService) GetCategoryName(ctx context.Context, categoryID int) (string, error) {
	return s.repo.GetCategoryName(ctx, categoryID)
}

func (s *ThreadService) Create(ctx context.Context, thread *Thread) (int, error) {
	if err := validateThreadContent(thread.Title, thread.Body); err != nil {
		return 0, err
	}
	// No need to Validate AuthorID since user was already authenticated
	// If CategoryID or AuthorID doesn't exist, there will be a repo level foreign key constraint error CategoryID exists

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

func (s *ThreadService) DisplayThreadPage(ctx context.Context, threadID int, username string, renderer *render.Renderer) (*ThreadData, error) {
	if threadID <= 0 {
		return nil, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}

	data, err := s.repo.DisplayThreadPage(ctx, threadID, username)

	var errFormat error
	data.Thread.DateCreated, errFormat = config.FormatTime(data.Thread.DateCreated)
	if errFormat != nil {
		return nil, errFormat
	}

	data.Thread.BodyHTML, err = renderer.ParseHTML(data.Thread.Body)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage body ParseHTML err: %w", err)
	}

	data.Thread.TitleHTML, err = renderer.ParseHTML(data.Thread.Title)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage title ParseHTML err: %w", err)
	}

	for i, comment := range data.Comments {
		data.Comments[i].DateCreated, errFormat = config.FormatTime(comment.DateCreated)
		if errFormat != nil {
			return nil, errFormat
		}

		data.Comments[i].BodyHTML, err = renderer.ParseHTML(data.Comments[i].Body)
		if err != nil {
			return nil, fmt.Errorf("DisplayThreadPage comment body ParseHTML err: %w", err)
		}
	}

	return data, err
}
