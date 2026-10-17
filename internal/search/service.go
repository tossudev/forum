package search

import (
	"context"

	"forum/internal/config"
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/errs"
	"forum/internal/pagination"
)

type SearchService struct {
	repo *SearchRepository
}

func NewService(r *SearchRepository) *SearchService {
	return &SearchService{repo: r}
}

func (s *SearchService) Search(ctx context.Context, input string, pagination pagination.Pagination) ([]thread.Thread, []comment.Comment, error) {
	threads, comments, err := s.repo.Search(ctx, input, pagination)
	if err != nil {
		return nil, nil, err
	}

	if len(threads) == 0 && len(comments) == 0 {
		return nil, nil, errs.ErrSearchNoResults
	}

	var errFormat error
	for i, thread := range threads {
		threads[i].DateCreated, errFormat = config.FormatTime(thread.DateCreated)
		if errFormat != nil {
			return nil, nil, errFormat
		}
	}

	for j, comment := range comments {
		comments[j].DateCreated, errFormat = config.FormatTime(comment.DateCreated)
		if errFormat != nil {
			return nil, nil, errFormat
		}
	}

	return threads, comments, nil
}
