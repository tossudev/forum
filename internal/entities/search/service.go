package search

import (
	"fmt"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/entities/comments"
	"forum/internal/entities/threads"
)

type SearchService struct {
	//here goes maybe the service layers of comment and thread, maybe something else too?	
	CommentService 	*comments.CommentService
	ThreadService 	*threads.ThreadService
	
}

func NewService()


func (ss *SearchService) Search(input string) ([]Comment, []Thread) {
	return nil, nil 
}

