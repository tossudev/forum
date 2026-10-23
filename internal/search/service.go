package search

import (
	"context"
	"fmt"

	//"forum/internal/entities/comment"
	//"forum/internal/entities/thread"
	"forum/internal/errs"
	"forum/internal/pagination"
)

type SearchService struct {
	repo *SearchRepository
}

func NewService(r *SearchRepository) *SearchService {
	return &SearchService{repo: r}
}

func (s *SearchService) Search(ctx context.Context, searchQuery string, tab string, pagination pagination.Pagination) (PaginatedSearchResult, error) {

	if err := pagination.Validate(); err != nil {
		return PaginatedSearchResult{}, fmt.Errorf("%w: invalid pagination parameters", errs.ErrInvalidUserInput)
	}

	searchResult, err := s.repo.Search(ctx, searchQuery, tab, pagination)
	if err != nil {
		return PaginatedSearchResult{}, err
	}

	if len(searchResult.Threads) == 0 && len(searchResult.Comments) == 0 {
		return PaginatedSearchResult{}, errs.ErrSearchNoResults
	}

	return searchResult, nil
}
