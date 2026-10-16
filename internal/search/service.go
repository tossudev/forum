package search

import (
	"context"
	"time"

	"forum/internal/config"
	"forum/internal/pagination"
)

/*type SearchService struct {
	//here goes maybe the service layers of comment and thread, maybe something else too?
	CommentService 	*comments.CommentService
	ThreadService 	*threads.ThreadService

}*/

type SearchService struct {
	repo *SearchRepository
}

func NewService(r *SearchRepository) *SearchService {
	return &SearchService{repo: r}
}

func (s *SearchService) Search(ctx context.Context, input string, pagination pagination.Pagination) (*SearchResult, error) {
	data, err := s.repo.Search(ctx, input, pagination)

	for i, thread := range data.Threads {
		timeThread, errTime := time.Parse(config.TimeFormat, thread.DateCreated)
		if errTime != nil {
			return nil, errTime
		}
		timeFormatted := timeThread.Format(config.TimeFormatFrontEnd)
		data.Threads[i].DateCreated = timeFormatted
	}

	for j, comment := range data.Comments {
		timeComment, errTime := time.Parse(config.TimeFormat, comment.DateCreated)
		if errTime != nil {
			return nil, errTime
		}
		timeCommentFormatted := timeComment.Format(config.TimeFormatFrontEnd)
		data.Comments[j].DateCreated = timeCommentFormatted
	}

	return data, err
}
