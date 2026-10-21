package thread

import (
	"forum/internal/pagination"
	"html/template"
)

type Thread struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	DateCreated string `json:"date_created"`
	AuthorID    int    `json:"author_id"`
	CategoryID  int    `json:"category_id"`
	ReplyCount  int
	TitleHTML   template.HTML
	BodyHTML    template.HTML
}

type ThreadsPage struct {
	CategoryID int
	Path       string
	Threads    []Thread
	PageNav    pagination.PaginationNav
	Username   string
}

type CreateThreadsPage struct {
	CategoryID   int
	CategoryName string
	InputTitle   string
	InputBody    string
	TitleErr     string
	BodyErr      string
	OtherErr     string
	Username     string
}

// Data for thread page with coments
type ThreadData struct {
	Path         string
	Thread       Thread
	AuthorName   string
	AuthorAvatar string
	LikeData     LikeData // LikeData for the thread post only
	Comments     []CommentData
	Username     string
}

type LikeData struct {
	NumLikes           int
	NumDislikes        int
	LoggedUserLiked    bool
	LoggedUserDisliked bool
}

type CommentData struct {
	ID          int
	Body        string
	BodyHTML    template.HTML
	DateCreated string
	ThreadID    int
	AuthorName  string
	LikeData    LikeData
}
