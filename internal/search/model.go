package search

import (
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/pagination"
)

type PaginatedSearchResult struct {
	Threads       []thread.Thread
	ThreadsCount  int
	Comments      []comment.Comment
	CommentsCount int
}

type SearchResultPage struct {
	ThreadPath     string
	CommentPath    string
	InputQuery     string
	Threads        []thread.Thread
	Comments       []comment.Comment
	ThreadPageNav  pagination.PaginationNav
	CommentPageNav pagination.PaginationNav
	Username       string
	ResponseMsg    string
}
