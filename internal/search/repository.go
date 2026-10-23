package search

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/config"
	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/middleware"
	"forum/internal/pagination"
)

type SearchRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) Search(ctx context.Context, input string, pagination pagination.Pagination) ([]thread.Thread, int, []comment.Comment, int, error) {
	queryInput := "\"" + input + "\""
	queryThread := `SELECT COUNT(*) OVER() AS total_count, rowid, title, body, date_created
	FROM threads_fts
	WHERE threads_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	var threadsCount int

	rowsThread, err := r.db.QueryContext(ctx, queryThread, queryInput, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, nil, 0, fmt.Errorf("Thread Search: %w", err)
	}

	var threads []thread.Thread
	var errFormat error
	for rowsThread.Next() {
		var thread thread.Thread
		if err := rowsThread.Scan(&threadsCount, &thread.ID, &thread.Title, &thread.Body, &thread.DateCreated); err != nil {
			return nil, 0, nil, 0, fmt.Errorf("Thread Search: %w", err)
		}

		thread.DateCreated, errFormat = config.FormatTime(thread.DateCreated, middleware.GetTimezone(ctx))
		if errFormat != nil {
			return nil, 0, nil, 0, errFormat
		}

		threads = append(threads, thread)
	}
	if err = rowsThread.Err(); err != nil {
		return nil, 0, nil, 0, fmt.Errorf("Thread Search: %w", err)
	}

	queryComment := `SELECT COUNT(*) OVER() AS total_count, body, date_created, thread_id
	FROM comments_fts
	WHERE comments_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	var commentsCount int

	rowsComment, err := r.db.QueryContext(ctx, queryComment, queryInput, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, nil, 0, fmt.Errorf("Comment Search: %w", err)
	}

	var comments []comment.Comment
	for rowsComment.Next() {
		var comment comment.Comment
		if err := rowsComment.Scan(&commentsCount, &comment.Body, &comment.DateCreated, &comment.ThreadID); err != nil {
			return nil, 0, nil, 0, fmt.Errorf("Comment Search: %w", err)
		}

		comment.DateCreated, errFormat = config.FormatTime(comment.DateCreated, middleware.GetTimezone(ctx))
		if errFormat != nil {
			return nil, 0, nil, 0, errFormat
		}

		comments = append(comments, comment)
	}
	if err = rowsComment.Err(); err != nil {
		return nil, 0, nil, 0, fmt.Errorf("Comment Search: %w", err)
	}

	return threads, threadsCount, comments, commentsCount, nil
}
