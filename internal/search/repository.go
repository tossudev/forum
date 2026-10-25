package search

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/entities/comment"
	"forum/internal/entities/thread"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/timeformat"
)

type SearchRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) Search(ctx context.Context, searchQuery string, tab string, paginationData pagination.Pagination) (PaginatedSearchResult, error) {

	//add double quotes around searchQuery so that single quotes (ex apostrophes) don't cause errors
	queryInput := "\"" + searchQuery + "\""

	var threads []thread.Thread
	var threadsCount int
	var comments []comment.Comment
	var commentsCount int
	var err error

	//if tab is "threads", search only threads
	if tab == "threads" {

		//get threads matching search term on requested page and total threads count
		threads, threadsCount, err = r.searchThreads(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search thread: %w", err)
		}
		return PaginatedSearchResult{Threads: threads, ThreadsCount: threadsCount}, nil

		//if tab is "comments", search only comments
	} else if tab == "comments" {

		//get comments matching search term on requested page and total comments count
		comments, commentsCount, err = r.searchComments(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search comments: %w", err)
		}
		return PaginatedSearchResult{Comments: comments, CommentsCount: commentsCount}, nil

		//if no tab (or if tab is anything other than "threads" or "comments", search both
	} else {

		//get threads matching search term on requested page and total threads count
		threads, threadsCount, err = r.searchThreads(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search threads: %w", err)
		}
		//get comments matching search term on requested page and total comments count
		comments, commentsCount, err = r.searchComments(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search comments: %w", err)
		}
	}

	searchResult := PaginatedSearchResult{
		Threads:       threads,
		ThreadsCount:  threadsCount,
		Comments:      comments,
		CommentsCount: commentsCount,
	}
	return searchResult, nil
}

func (r *SearchRepository) searchThreads(ctx context.Context, queryInput string, paginationData pagination.Pagination) ([]thread.Thread, int, error) {

	var threads []thread.Thread
	var threadsCount int

	//get count of threads that match search term
	queryCount := `SELECT COUNT(*) FROM threads_fts WHERE threads_fts MATCH ?;`
	err := r.db.QueryRowContext(ctx, queryCount, queryInput).Scan(&threadsCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Search threads: %w", err)
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, threadsCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Search Thread: %w", errs.ErrInvalidUserInput)
	}

	//get threads
	queryThread := `SELECT rowid, title, body, date_created
	FROM threads_fts
	WHERE threads_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsThread, err := r.db.QueryContext(ctx, queryThread, queryInput, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("Thread Search: %w", err)
	}

	var errFormat error
	for rowsThread.Next() {
		var thread thread.Thread
		if err := rowsThread.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated); err != nil {
			return nil, 0, fmt.Errorf("Thread Search: %w", err)
		}

		thread.DateCreated, errFormat = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
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

func (r *SearchRepository) searchComments(ctx context.Context, queryInput string, paginationData pagination.Pagination) ([]comment.Comment, int, error) {

	var comments []comment.Comment
	var commentsCount int

	//get count of comments that match search term
	queryCount := `SELECT COUNT(*) FROM comments_fts WHERE comments_fts MATCH ?;`
	err := r.db.QueryRowContext(ctx, queryCount, queryInput).Scan(&commentsCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Search comments: %w", err)
	}

	//check if requested page is larger than possible pages based on number of results & page size
	//if so, set page to last possible page
	err = pagination.CheckPageTooLarge(paginationData, commentsCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Search comments: %w", errs.ErrInvalidUserInput)
	}

	//get comments
	queryComment := `SELECT body, date_created, thread_id
	FROM comments_fts
	WHERE comments_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsComment, err := r.db.QueryContext(ctx, queryComment, queryInput, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("Comment Search: %w", err)
	}

	var errFormat error
	for rowsComment.Next() {
		var comment comment.Comment
		if err := rowsComment.Scan(&comment.Body, &comment.DateCreated, &comment.ThreadID); err != nil {
			return nil, 0, fmt.Errorf("Comment Search: %w", err)
		}

		comment.DateCreated, errFormat = timeformat.FormatTime(comment.DateCreated, timeformat.GetTimezone(ctx))
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
