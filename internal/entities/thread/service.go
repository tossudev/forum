package thread

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/validate"
)

type ThreadService struct {
	repo *ThreadRepository
}

func NewService(r *ThreadRepository) *ThreadService {
	return &ThreadService{repo: r}
}

// GetByCategory gets threads (paginated) by category ID.
// Returns the threads, total found, the category name, possibly updated page number, and a potential error.
func (s *ThreadService) GetByCategory(ctx context.Context, id int, pagination pagination.Pagination) ([]Thread, int, string, int, error) {
	if id <= 0 {
		return nil, 0, "", 0, fmt.Errorf("%w: invalid category id", errs.ErrInvalidUserInput)
	}
	// This validates that the category exists and also fetches the category name
	categoryName, err := s.repo.GetCategoryName(ctx, id)
	if err != nil {
		return nil, 0, "", 0, err
	}

	if err := pagination.Validate(); err != nil {
		return nil, 0, "", 0, fmt.Errorf("%w: invalid pagination parameters", errs.ErrInvalidUserInput)
	}

	threads, totalCount, page, err := s.repo.GetByCategory(ctx, id, pagination)
	if err != nil {
		return nil, 0, "", 0, err
	}

	for i, thread := range threads {
		threads[i].TitleHTML, err = validate.ParseHTML(thread.Title)
		if err != nil {
			threads[i].DisplayError = err
			slog.Error("thread title validation failed", "err:", err, "thread_id:", threads[i].ID, "title:", threads[i].Title)
		}
	}

	return threads, totalCount, categoryName, page, nil
}

func (s *ThreadService) GetCategoryName(ctx context.Context, categoryID int) (string, error) {
	return s.repo.GetCategoryName(ctx, categoryID)
}

func (s *ThreadService) Create(ctx context.Context, thread *Thread) (int, error) {
	if err := validateThreadContent(thread.Title, thread.Body); err != nil {
		return 0, err
	}

	if _, err := validate.ParseHTML(thread.Title); err != nil {
		return 0, err
	}

	if _, err := validate.ParseHTML(thread.Body); err != nil {
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

// DisplayThreadPage gets the data needed for displaying a single thread page (including comments and associated likes).
// Also returns an int of the total comments count for pagination.
func (s *ThreadService) DisplayThreadPage(ctx context.Context, threadID int, username string, paginationData pagination.Pagination) (*ThreadDisplay, []CommentDisplay, int, error) {
	if threadID <= 0 {
		return nil, nil, 0, errs.ErrNotFound
	}

	// Check if thread exists
	if _, err := s.repo.GetByID(ctx, threadID); err != nil {
		return nil, nil, 0, err
	}

	threadDisplay, commentsDisplay, totalCount, err := s.repo.DisplayThreadPage(ctx, threadID, username, paginationData)
	if err != nil {
		return nil, nil, 0, err
	}

	//TODO: Could refactor to move the HTML sanitization into the scan loops at the repo layer to avoid double looping

	threadDisplay.Thread.TitleHTML, err = validate.ParseHTML(threadDisplay.Thread.Title)
	if err != nil {
		threadDisplay.DisplayError = err
		slog.Error("thread title validation failed", "err:", err, "thread_id:", threadDisplay.Thread.ID, "title:", threadDisplay.Thread.Title)
	}

	threadDisplay.Thread.BodyHTML, err = validate.ParseHTML(threadDisplay.Thread.Body)
	if err != nil {
		threadDisplay.DisplayError = err
		slog.Error("thread body validation failed", "err:", err, "thread_id:", threadDisplay.Thread.ID, "body:", threadDisplay.Thread.Body)
	}

	for i, c := range commentsDisplay {
		commentsDisplay[i].Comment.BodyHTML, err = validate.ParseHTML(c.Comment.Body)
		if err != nil {
			commentsDisplay[i].DisplayError = err
			slog.Error("comment validation failed", "err:", err, "comment_id:", commentsDisplay[i].Comment.ID, "body:", commentsDisplay[i].Comment.Body)
		}
	}

	//TODO: Check if this should still return nil if there is a display error (as above)
	return threadDisplay, commentsDisplay, totalCount, nil
}
