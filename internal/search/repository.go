package search

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/config"
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

// TODO: send in tab, switch different tab possibilities, create helper to run queries
// Need to send back max number of pages for search
func (r *SearchRepository) Search(ctx context.Context, input string, tab string, paginationData pagination.Pagination) ([]thread.Thread, int, []comment.Comment, int, error) {

	//add double quotes around input so that single quotes (ex apostrophes) don't cause errors
	queryInput := "\"" + input + "\""

	//check if page number is beyond possible pages
	//if so, make page number equal last page
	page, err := r.checkPage(ctx, queryInput, tab, paginationData)
	if err != nil {
		return nil, 0, nil, 0, fmt.Errorf("Search: %w,", err)
	}
	paginationData.Page = page

	var threads []thread.Thread
	var threadsCount int
	var comments []comment.Comment
	var commentsCount int

	//if tab is "threads", search only threads
	if tab == "threads" {
		threads, threadsCount, err = r.searchThreads(ctx, queryInput, paginationData)
		if err != nil {
			return nil, 0, nil, 0, fmt.Errorf("Search threads: %w", err)
		}
		return threads, threadsCount, nil, 0, nil

		//if tab is "comments", search only comments
	} else if tab == "comments" {
		comments, commentsCount, err = r.searchComments(ctx, queryInput, paginationData)
		if err != nil {
			return nil, 0, nil, 0, fmt.Errorf("Search comments: %w", err)
		}
		return nil, 0, comments, commentsCount, nil

		//if no tab (or if tab is anything other than "threads" or "comments", search both
	} else {
		threads, threadsCount, err = r.searchThreads(ctx, queryInput, paginationData)
		if err != nil {
			return nil, 0, nil, 0, fmt.Errorf("Search threads: %w", err)
		}

		comments, commentsCount, err = r.searchComments(ctx, queryInput, paginationData)
		if err != nil {
			return nil, 0, nil, 0, fmt.Errorf("Search comments: %w", err)
		}
	}

	return threads, threadsCount, comments, commentsCount, nil
}

func (r *SearchRepository) checkPage(ctx context.Context, input string, tab string, paginationData pagination.Pagination) (int, error) {
	if tab != "threads" && tab != "comments" {
		return paginationData.Page, nil
	}

	var query string
	var count int
	if tab == "threads" {
		query = `SELECT COUNT(*) OVER() as total_count
		FROM threads_fts
		WHERE threads_fts MATCH ?;`
	}

	if tab == "comments" {
		query = `SELECT COUNT(*) OVER() as total_count
		FROM comments_fts
		WHERE comments_fts MATCH ?;`
	}

	err := r.db.QueryRowContext(ctx, query, input).Scan(&count)
	if err != nil {
		return 0, err
	}

	if paginationData.Page > pagination.CalcTotalPages(count, paginationData.PageSize) {
		return pagination.CalcTotalPages(count, paginationData.PageSize), nil
	}

	return paginationData.Page, nil
}

func (r *SearchRepository) searchThreads(ctx context.Context, input string, paginationData pagination.Pagination) ([]thread.Thread, int, error) {

	//get threads
	queryThread := `SELECT COUNT(*) OVER() AS total_count, rowid, title, body, date_created
	FROM threads_fts
	WHERE threads_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	var threadsCount int

	rowsThread, err := r.db.QueryContext(ctx, queryThread, input, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("Thread Search: %w", err)
	}

	var threads []thread.Thread
	var errFormat error
	for rowsThread.Next() {
		var thread thread.Thread
		if err := rowsThread.Scan(&threadsCount, &thread.ID, &thread.Title, &thread.Body, &thread.DateCreated); err != nil {
			return nil, 0, fmt.Errorf("Thread Search: %w", err)
		}

		thread.DateCreated, errFormat = config.FormatTime(thread.DateCreated)
		if errFormat != nil {
			return nil, 0, errFormat
		}

		threads = append(threads, thread)
	}
	if err = rowsThread.Err(); err != nil {
		return nil, 0, fmt.Errorf("Thread Search: %w", err)
	}

	return threads, threadsCount, nil
}

func (r *SearchRepository) searchComments(ctx context.Context, input string, paginationData pagination.Pagination) ([]comment.Comment, int, error) {

	//get comments
	queryComment := `SELECT COUNT(*) OVER() AS total_count, body, date_created, thread_id
	FROM comments_fts
	WHERE comments_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsComment, err := r.db.QueryContext(ctx, queryComment, input, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("Comment Search: %w", err)
	}

	var comments []comment.Comment
	var commentsCount int
	var errFormat error
	for rowsComment.Next() {
		var comment comment.Comment
		if err := rowsComment.Scan(&commentsCount, &comment.Body, &comment.DateCreated, &comment.ThreadID); err != nil {
			return nil, 0, fmt.Errorf("Comment Search: %w", err)
		}

		comment.DateCreated, errFormat = config.FormatTime(comment.DateCreated)
		if errFormat != nil {
			return nil, 0, errFormat
		}

		comments = append(comments, comment)
	}

	if err = rowsComment.Err(); err != nil {
		return nil, 0, fmt.Errorf("Comment Search: %w", err)
	}

	return comments, commentsCount, nil
}
