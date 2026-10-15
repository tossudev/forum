package search

import (
	"context"

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
	//TODO: convert time format if it is displayed on frontent

	return s.repo.Search(ctx, input, pagination)
}
