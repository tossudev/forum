package search

import (
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
)

type SearchResult struct {
	Path       string
	InputQuery string
	Threads    []thread.Thread
	Comments   []comment.Comment
	Page       int
	Username   string
}
