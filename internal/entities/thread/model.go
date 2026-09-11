package thread

type Thread struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	DateCreated string `json:"date_created"`
	AuthorID    int    `json:"author_id"`
	CategoryID  int    `json:"category_id"`
}

type ThreadRequest struct {
	Title       *string `json:"title"`
	Body        *string `json:"body"`
	DateCreated *string `json:"date_created"`
	AuthorID    *int    `json:"author_id"`
	CategoryID  *int    `json:"category_id"`
}
