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

//func (s *SearchService) Search(ctx context.Context, searchQuery string, tab string, pagination pagination.Pagination) ([]thread.Thread, int, []comment.Comment, int, int, error) {

func (s *SearchService) Search(ctx context.Context, searchQuery string, tab string, pagination pagination.Pagination) (PaginatedSearchResult, error) {

	if err := pagination.Validate(); err != nil {
		//return nil, 0, nil, 0, 0, fmt.Errorf("%w: invalid pagination parameters", errs.ErrInvalidUserInput)
		return PaginatedSearchResult{}, fmt.Errorf("%w: invalid pagination parameters", errs.ErrInvalidUserInput)
	}

	//threads, threadsCount, comments, commentsCount, page, err := s.repo.Search(ctx, searchQuery, tab, pagination)
	searchResult, err := s.repo.Search(ctx, searchQuery, tab, pagination)
	if err != nil {
		//return nil, 0, nil, 0, 0, err
		return PaginatedSearchResult{}, err
	}

	if len(searchResult.Threads) == 0 && len(searchResult.Comments) == 0 {
		//return nil, 0, nil, 0, 0, errs.ErrSearchNoResults
		return PaginatedSearchResult{}, errs.ErrSearchNoResults
	}

	//return threads, threadsCount, comments, commentsCount, page, nil
	return searchResult, nil
}
