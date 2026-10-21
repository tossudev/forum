package search

import (
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/pagination"
)

type SearchResult struct {
	InputQuery     string
	Threads        []thread.Thread
	Comments       []comment.Comment
	ThreadPageNav  pagination.PaginationNav
	CommentPageNav pagination.PaginationNav
	Username       string
	ResponseMsg    string
}
