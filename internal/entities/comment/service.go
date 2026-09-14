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

func (s *CommentService) Create(ctx context.Context, req *Comment) (Comment, error) {

	isoTime := time.Now().UTC()
	isoString := isoTime.Format(time.RFC3339)

	req.DateCreated = isoString

	newComment, err := s.repo.Create(ctx, req)
	if err != nil {
		return Comment{}, err
	}

	return newComment, nil
}

func (s *CommentService) GetByID(ctx context.Context, id int) (Comment, error) {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil

}

func (s *CommentService) GetByThread(ctx context.Context, id int) ([]Comment, error) {
	return nil, nil
}
