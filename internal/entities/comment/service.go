package comment

import (
	"context"
	"fmt"
	"time"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
)

type CommentService struct {
	repo *CommentRepository
}

func NewService(r *CommentRepository) *CommentService {
	return &CommentService{repo: r}
}

func (s *CommentService) Create(ctx context.Context, req *Comment) (Comment, error) {

	isoTime := time.Now().UTC()
	isoString := isoTime.Format(config.TimeFormat)

	req.DateCreated = isoString

	newComment, err := s.repo.Create(ctx, req)
	if err != nil {
		return Comment{}, err
	}

	return newComment, nil
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
	//TODO: Logic for the service layer --> who gets to delete a comment (user, admin etc.)
	if commentID <= 0 {
		return fmt.Errorf("%w: invalid comment id", errs.ErrInvalidUserInput)
	}

	return s.repo.Delete(ctx, commentID, userID)
}

func (s *CommentService) Filter(ctx context.Context, input string) error {
	//TODO
	return nil
}
