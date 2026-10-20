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

type ThreadData struct {
	Path         string
	Thread       Thread
	AuthorName   string
	AuthorAvatar string
	NumLikes     int
	LikeUsers    []string
	NumDislikes  int
	DislikeUsers []string
	Comments     []CommentData
	Username     string
	Avatar       string
}

type CommentData struct {
	ID                 int
	Body               string
	BodyHTML           template.HTML
	DateCreated        string
	ThreadID           int
	AuthorName         string
	NumLikes           int
	LikeUsers          []string
	NumDislikes        int
	DislikeUsers       []string
	LoggedUserLiked    bool
	LoggedUserDisliked bool
	Avatar             string
}
