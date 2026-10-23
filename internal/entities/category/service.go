package category

import (
	"context"

	"forum/internal/pagination"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewService(r *CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) GetAll(ctx context.Context, pagination pagination.Pagination) ([]Category, error) {
	return s.repo.GetAll(ctx, pagination)
}

func (s *CategoryService) Create(ctx context.Context, category *Category) error {
	return s.repo.Create(ctx, category)
}
