package search

import (
	"context"
	"fmt"

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

func (s *SearchService) Search(ctx context.Context, input string, pagination pagination.Pagination) ([]thread.Thread, int, []comment.Comment, int, error) {

	if err := pagination.Validate(); err != nil {
		return nil, 0, nil, 0, fmt.Errorf("%w: invalid pagination parameters", errs.ErrInvalidUserInput)
	}

	threads, threadsCount, comments, commentsCount, err := s.repo.Search(ctx, input, pagination)
	if err != nil {
		return nil, 0, nil, 0, err
	}

	if len(threads) == 0 && len(comments) == 0 {
		return nil, 0, nil, 0, errs.ErrSearchNoResults
	}

	var errFormat error
	for i, thread := range threads {
		threads[i].DateCreated, errFormat = config.FormatTime(thread.DateCreated)
		if errFormat != nil {
			return nil, 0, nil, 0, errFormat
		}
	}

	for j, comment := range comments {
		comments[j].DateCreated, errFormat = config.FormatTime(comment.DateCreated)
		if errFormat != nil {
			return nil, 0, nil, 0, errFormat
		}
	}

	return threads, threadsCount, comments, commentsCount, nil
}
