package thread

import (
	"forum/internal/models"
	"forum/internal/pagination"
	"html/template"
)

type Thread struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	DateCreated  string `json:"date_created"`
	AuthorID     int    `json:"author_id"`
	CategoryID   int    `json:"category_id"`
	ReplyCount   int
	TitleHTML    template.HTML
	BodyHTML     template.HTML
	DisplayError error // TODO: Move this out and don't add threads to db that don't pass validation?
	// Now for displaying a thread page with comments, the display error is in ThreadDisplay,
	// but for viewing threads in a category, it's directly in the thread model
}

// ThreadsPage contains the page data for viewing the threads in a category
type ThreadsPage struct {
	CategoryID  int
	Breadcrumbs []models.Breadcrumb
	Threads     []Thread
	PageNav     pagination.PaginationNav
	Username    string
}

// CreateThreadsPage contains the page data for create a thread page
type CreateThreadsPage struct {
	CategoryID   int
	CategoryName string
	InputTitle   string
	InputBody    string
	TitleErr     string
	BodyErr      string
	OtherErr     string
	Username     string
	Path         string
	Threads      []Thread
	Page         int
}

// ThreadCommentsPage contains page data for viewing the comments within a thread
type ThreadCommentsPage struct {
	Breadcrumbs   []models.Breadcrumb
	PageNav       pagination.PaginationNav
	Username      string
	ThreadDisplay ThreadDisplay
	Comments      []CommentDisplay
}

// ThreadDisplay contains data needed to display a thread box on the UI
type ThreadDisplay struct {
	Thread       Thread
	AuthorName   string
	AuthorAvatar string
	LikeData     LikeData
	DisplayError error
}

// CommentDisplay contains data needed to display a comment on the UI
type CommentDisplay struct {
	Comment      Comment
	AuthorName   string
	AuthorAvatar string
	LikeData     LikeData
	DisplayError error
}

// LikeData contains data needed to display likes on the UI
type LikeData struct {
	NumLikes           int
	NumDislikes        int
	LoggedUserLiked    bool
	LoggedUserDisliked bool
	DisplayError       error
}

// Temporary mock of the same Comment model from comment package - to be moved to shared models
type Comment struct {
	ID          int    `json:"id"`
	Body        string `json:"body"`
	DateCreated string `json:"date_created"`
	ThreadID    int    `json:"thread_id"`
	AuthorID    int    `json:"author_id"`
	BodyHTML    template.HTML
}
