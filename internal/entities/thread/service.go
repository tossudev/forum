package thread

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"forum/internal/errs"
	"forum/internal/models"
	"forum/internal/pagination"
	"forum/internal/timeformat"
	"forum/internal/validate"
)

type ThreadService struct {
	repo *ThreadRepository
}

func NewService(r *ThreadRepository) *ThreadService {
	return &ThreadService{repo: r}
}

// GetByCategory gets threads (paginated) by category ID.
// Returns the threads, total found, the category name, and a potential error.
func (s *ThreadService) GetByCategory(ctx context.Context, id int, pagination pagination.Pagination) ([]Thread, int, string, error) {
	// This validates that the category exists and also fetches the category name
	categoryName, err := s.repo.GetCategoryName(ctx, id)
	if err != nil {
		return nil, 0, "", err
	}

	if err := pagination.Validate(); err != nil {
		return nil, 0, "", errs.ErrInvalidUserInput
	}

	threads, totalCount, err := s.repo.GetByCategory(ctx, id, pagination)
	if err != nil {
		return nil, 0, "", err
	}

	for i, thread := range threads {
		threads[i].TitleHTML, err = validate.ParseHTML(thread.Title)
		if err != nil {
			threads[i].DisplayError = err
			slog.Error("thread title validation failed", "err:", err, "thread_id:", threads[i].ID, "title:", threads[i].Title)
		}
	}

	return threads, totalCount, categoryName, nil
}

func (s *ThreadService) GetCategoryName(ctx context.Context, categoryID int) (string, error) {
	return s.repo.GetCategoryName(ctx, categoryID)
}

func (s *ThreadService) Create(ctx context.Context, thread *Thread) (int, error) {
	if err := validateThreadContent(thread.Title, thread.Body); err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	if _, err := validate.ParseHTML(thread.Title); err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	if _, err := validate.ParseHTML(thread.Body); err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	now := time.Now().UTC().Format(timeformat.TimeFormat)
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
	// Check if thread exists
	if _, err := s.repo.GetByID(ctx, threadID); err != nil {
		return nil, nil, 0, err
	}

	threadDisplay, commentsDisplay, totalCount, err := s.repo.DisplayThreadPage(ctx, threadID, username, paginationData)
	if err != nil {
		return nil, nil, totalCount, err // totalCount needed when err is ErrPageOverTotal for redirection to last page
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

func (s *ThreadService) CollectPageData(ctx context.Context, params ThreadPageParams) (ThreadCommentsPage, int, error) {

	threadDisplay, commentsDisplay, totalCount, err := s.DisplayThreadPage(ctx, params.ThreadID, params.Username, params.PaginationData)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return ThreadCommentsPage{}, totalCount, errs.ErrNotFound
		}
		return ThreadCommentsPage{}, totalCount, err
	}

	catID := threadDisplay.Thread.CategoryID
	categoryName, err := s.GetCategoryName(ctx, catID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return ThreadCommentsPage{}, totalCount, errs.ErrNotFound
		}
		return ThreadCommentsPage{}, totalCount, err
	}

	breadcrumbs := []models.Breadcrumb{
		models.Breadcrumb{
			Name: "Home",
			URL:  "/",
		},
		models.Breadcrumb{
			Name: categoryName,
			URL:  "/categories/" + strconv.Itoa(catID),
		},
		models.Breadcrumb{
			Name: threadDisplay.Thread.Title,
			URL:  "/threads/" + strconv.Itoa(params.ThreadID),
		},
	}
	totalPages := pagination.CalcTotalPages(totalCount, params.PaginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(params.PaginationData.Page, totalPages, 5) // The visible pages in navigation
	pageNav := pagination.PaginationNav{
		Path:       params.Path,
		Page:       params.PaginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   params.PaginationData.Page - 1,
		NextPage:   params.PaginationData.Page + 1,
	}
	pageData := ThreadCommentsPage{
		Breadcrumbs:   breadcrumbs,
		PageNav:       pageNav,
		Username:      params.Username,
		ThreadDisplay: *threadDisplay,
		Comments:      commentsDisplay,
	}
	return pageData, totalCount, nil
}
