package category

import (
	"context"
	"fmt"

	"forum/internal/errs"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewService(r *CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (*Category, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid thread id", errs.ErrInvalidUserInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, category *Category) error {
	return s.repo.Create(ctx, category)
}
