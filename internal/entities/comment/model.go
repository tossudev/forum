package comment

import (
	"html/template"
)

type Comment struct {
	ID          int    `json:"id"`
	Body        string `json:"body"`
	DateCreated string `json:"date_created"`
	ThreadID    int    `json:"thread_id"`
	AuthorID    int    `json:"author_id"`
	BodyHTML    template.HTML
}

type ErrorPage struct {
	Username string
	Message  string
}
