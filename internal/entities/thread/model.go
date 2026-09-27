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
	Username     string
}
