package comment

import ()

type CommentService struct {
	repo *CommentRepository
}

func NewService(r *CommentRepository) *CommentService {
	return &CommentService{repo: r}
}

func (s *CommentService) Create(comment Comment) error {
	return nil
}

func (s *CommentService) GetByThread(id int) error {
	return nil
}
