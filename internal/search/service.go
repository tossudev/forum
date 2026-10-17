package search

import (
	"context"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
)

type SearchService struct {
	repo *SearchRepository
}

func NewService(r *SearchRepository) *SearchService {
	return &SearchService{repo: r}
}

func (s *SearchService) Search(ctx context.Context, input string, pagination pagination.Pagination) (*SearchResult, error) {
	data, err := s.repo.Search(ctx, input, pagination)

	if len(data.Threads) == 0 && len(data.Comments) == 0 {
		return data, errs.ErrSearchNoResults
	}

	var errFormat error
	for i, thread := range data.Threads {
		data.Threads[i].DateCreated, errFormat = config.FormatTime(thread.DateCreated)
		if errFormat != nil {
			return nil, errFormat
		}
	}

	for j, comment := range data.Comments {
		data.Comments[j].DateCreated, errFormat = config.FormatTime(comment.DateCreated)
		if errFormat != nil {
			return nil, errFormat
		}
	}

	return data, err
}
