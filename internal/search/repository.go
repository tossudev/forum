package search

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/pagination"
)

type SearchRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) Search(ctx context.Context, input string, pagination pagination.Pagination) (*SearchResult, error) {
	//TODO: look at query parameters (ex: search partial words, etc)
	//TODO: add a note on search page for user that search works for title and body in threads, comments

	queryThread := `SELECT rowid, title, body, date_created
	FROM threads_fts
	WHERE threads_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsThread, err := r.db.QueryContext(ctx, queryThread, input, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, fmt.Errorf("Thread Search: %w", err)
	}

	var data SearchResult
	var threads []thread.Thread
	for rowsThread.Next() {
		var thread thread.Thread
		if err := rowsThread.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated); err != nil {
			return nil, fmt.Errorf("Thread Search: %w", err)
		}
		threads = append(threads, thread)
	}
	if err = rowsThread.Err(); err != nil {
		return nil, fmt.Errorf("Thread Search: %w", err)
	}

	data.Threads = threads

	queryComment := `SELECT body, date_created, thread_id
	FROM comments_fts
	WHERE comments_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsComment, err := r.db.QueryContext(ctx, queryComment, input, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, fmt.Errorf("Comment Search: %w", err)
	}

	var comments []comment.Comment
	for rowsComment.Next() {
		var comment comment.Comment
		if err := rowsComment.Scan(&comment.Body, &comment.DateCreated, &comment.ThreadID); err != nil {
			return nil, fmt.Errorf("Comment Search: %w", err)
		}
		comments = append(comments, comment)
	}
	if err = rowsComment.Err(); err != nil {
		return nil, fmt.Errorf("Comment Search: %w", err)
	}

	data.Comments = comments

	data.InputQuery = input

	return &data, nil
}
