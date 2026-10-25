package search

import (
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/models"
	"forum/internal/pagination"
)

type PaginatedSearchResult struct {
	Threads       []thread.Thread
	ThreadsCount  int
	Comments      []comment.Comment
	CommentsCount int
}

type SearchResultPage struct {
	ThreadBreadcrumbs  ThreadBreadcrumbs
	CommentBreadcrumbs CommentBreadcrumbs
	InputQuery         string
	Threads            []thread.Thread
	Comments           []comment.Comment
	ThreadPageNav      pagination.PaginationNav
	CommentPageNav     pagination.PaginationNav
	Username           string
	ResponseMsg        string
}

type ThreadBreadcrumbs struct {
	Breadcrumbs []models.Breadcrumb
}

type CommentBreadcrumbs struct {
	Breadcrumbs []models.Breadcrumb
}
