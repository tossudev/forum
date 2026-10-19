package thread

import (
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
	Page       int   // Current (wanted) page
	TotalPages int   // Last possible page
	PageList   []int // Visible pages in navigation
	PrevPage   int   // The page after the current page
	NextPage   int   // The page after the current page
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
}
