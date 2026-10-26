package comment

import (
	"context"
	"fmt"
	"time"

	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/timeformat"
	"forum/internal/validate"
)

type CommentService struct {
	repo *CommentRepository
}

func NewService(r *CommentRepository) *CommentService {
	return &CommentService{repo: r}
}

// Create creates a new comment, and returns the added comment, the total number of comments, and a potential error.
func (s *CommentService) Create(ctx context.Context, req *Comment) (Comment, int, error) {
	if _, err := validate.ParseHTML(req.Body); err != nil {
		return Comment{}, 0, errs.ErrInvalidUserInput
	}

	isoTime := time.Now().UTC()
	isoString := isoTime.Format(timeformat.TimeFormat)

	req.DateCreated = isoString

	newComment, totalComments, err := s.repo.Create(ctx, req)
	if err != nil {
		return Comment{}, 0, err
	}

	return newComment, totalComments, nil
}

func (s *CommentService) GetByID(ctx context.Context, id int) (Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) GetByThread(ctx context.Context, id int, p pagination.Pagination) ([]Comment, error) {

	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}
	comments, err := s.repo.GetByThread(ctx, id, p)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) Delete(ctx context.Context, commentID, userID int) error {
	if commentID <= 0 {
		return fmt.Errorf("%w: invalid comment id", errs.ErrInvalidUserInput)
	}

	return s.repo.Delete(ctx, commentID, userID)
}
