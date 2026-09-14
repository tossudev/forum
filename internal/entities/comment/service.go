package comment

import (
	"context"
	"time"
)

type CommentService struct {
	repo *CommentRepository
}

func NewService(r *CommentRepository) *CommentService {
	return &CommentService{repo: r}
}

func (s *CommentService) Create(ctx context.Context, c *Comment) error {

	isoTime := time.Now().UTC()
	isoString := isoTime.Format(time.RFC3339)

	c.DateCreated = isoString

	return nil
}

func (s *CommentService) GetByThread(ctx context.Context, id int) ([]Comment, error) {
	return nil, nil
}
