package search

import (
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
)

type SearchResult struct {
	ThreadPath        string
	CommentPath       string
	InputQuery        string
	Threads           []thread.Thread
	Comments          []comment.Comment
	ThreadPage        int
	ThreadTotalPages  int
	ThreadPageList    []int
	ThreadPrevPage    int
	ThreadNextPage    int
	CommentPage       int
	CommentTotalPages int
	CommentPageList   []int
	CommentPrevPage   int
	CommentNextPage   int
	Username          string
	ResponseMsg       string
}
