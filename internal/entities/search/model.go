package search

import (
	"forum/internal/entities/comment"
	"forum/internal/entities/threads"
)

type SearchResult struct {
	UserName 	string
	InputQuery	string
	Threads 	*[]thread.Thread
	Comments 	*[]comment.Comments
}
