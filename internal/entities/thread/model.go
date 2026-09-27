package thread

type Thread struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	DateCreated string `json:"date_created"`
	AuthorID    int    `json:"author_id"`
	CategoryID  int    `json:"category_id"`
}

type ThreadsPage struct {
	Path     string
	Threads  []Thread
	Page     int
	Username string
}

type ThreadData struct {
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
	ID           int
	Body         string
	DateCreated  string
	ThreadID     int
	AuthorName   string
	NumLikes     int
	LikeUsers    []string
	NumDislikes  int
	DislikeUsers []string
}
