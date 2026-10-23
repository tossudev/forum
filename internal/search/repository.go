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

func (r *SearchRepository) Search(ctx context.Context, searchQuery string, tab string, paginationData pagination.Pagination) (PaginatedSearchResult, error) {

	//add double quotes around searchQuery so that single quotes (ex apostrophes) don't cause errors
	queryInput := "\"" + searchQuery + "\""

	var threads []thread.Thread
	var threadsCount int
	var comments []comment.Comment
	var commentsCount int
	var page int

	//if tab is "threads", search only threads
	if tab == "threads" {

		//get count of threads that match search term
		queryCount := `SELECT COUNT(*) FROM threads_fts WHERE threads_fts MATCH ?;`
		err := r.db.QueryRowContext(ctx, queryCount, queryInput).Scan(&threadsCount)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search threads: %w", err)
		}

		//check if requested page is larger than last possible page based on number of results & page size
		//if so, set page to last possible page
		page = pagination.CheckPageTooLarge(paginationData, threadsCount)
		paginationData.Page = page

		//get threads matching search term on requested page
		threads, err = r.searchThreads(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search thread: %w", err)
		}
		return PaginatedSearchResult{Threads: threads, ThreadsCount: threadsCount, Page: page}, nil

		//if tab is "comments", search only comments
	} else if tab == "comments" {

		//get count of comments that match search term
		queryCount := `SELECT COUNT(*) FROM comments_fts WHERE comments_fts MATCH ?;`
		err := r.db.QueryRowContext(ctx, queryCount, queryInput).Scan(&commentsCount)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search comments: %w", err)
		}

		//check if requested page is larger than possible pages based on number of results & page size
		//if so, set page to last possible page
		page = pagination.CheckPageTooLarge(paginationData, commentsCount)
		paginationData.Page = page

		//get comments matching search term on requested page
		comments, err = r.searchComments(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search comments: %w", err)
		}
		return PaginatedSearchResult{Comments: comments, CommentsCount: commentsCount, Page: page}, nil

		//if no tab (or if tab is anything other than "threads" or "comments", search both
	} else {

		//get count of threads that match search term
		queryCountThreads := `SELECT COUNT(*) FROM threads_fts WHERE threads_fts MATCH ?;`
		err := r.db.QueryRowContext(ctx, queryCountThreads, queryInput).Scan(&threadsCount)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search threads: %w", err)
		}

		//get count of comments that match search term
		queryCountComments := `SELECT COUNT(*) FROM comments_fts WHERE comments_fts MATCH ?;`
		err = r.db.QueryRowContext(ctx, queryCountComments, queryInput).Scan(&commentsCount)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search comments: %w", err)
		}

		//if there are more comments than threads, get page based on number of results for comments
		//otherwise, get page based on number of results for threads
		if commentsCount > threadsCount {
			page = pagination.CheckPageTooLarge(paginationData, threadsCount)
		} else {
			page = pagination.CheckPageTooLarge(paginationData, commentsCount)
		}
		paginationData.Page = page

		//get threads matching search term on requested page
		threads, err = r.searchThreads(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search threads: %w", err)
		}

		//get comments matching search term on requested page
		comments, err = r.searchComments(ctx, queryInput, paginationData)
		if err != nil {
			return PaginatedSearchResult{}, fmt.Errorf("Search comments: %w", err)
		}
	}

	searchResult := PaginatedSearchResult{
		Threads:       threads,
		ThreadsCount:  threadsCount,
		Comments:      comments,
		CommentsCount: commentsCount,
		Page:          page,
	}
	return searchResult, nil
}

func (r *SearchRepository) searchThreads(ctx context.Context, searchQuery string, paginationData pagination.Pagination) ([]thread.Thread, error) {

	//get threads
	queryThread := `SELECT rowid, title, body, date_created
	FROM threads_fts
	WHERE threads_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsThread, err := r.db.QueryContext(ctx, queryThread, searchQuery, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, fmt.Errorf("Thread Search: %w", err)
	}

	var threads []thread.Thread
	var errFormat error
	for rowsThread.Next() {
		var thread thread.Thread
		if err := rowsThread.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated); err != nil {
			return nil, fmt.Errorf("Thread Search: %w", err)
		}

		thread.DateCreated, errFormat = config.FormatTime(thread.DateCreated, middleware.GetTimezone(ctx))
		if errFormat != nil {
			return nil, errFormat
		}

		threads = append(threads, thread)
	}
	if err = rowsThread.Err(); err != nil {
		return nil, fmt.Errorf("Thread Search: %w", err)
	}

	return threads, nil
}

func (r *SearchRepository) searchComments(ctx context.Context, searchQuery string, paginationData pagination.Pagination) ([]comment.Comment, error) {

	//get comments
	queryComment := `SELECT body, date_created, thread_id
	FROM comments_fts
	WHERE comments_fts MATCH ?
	ORDER BY rank
	LIMIT ? OFFSET ?;`

	rowsComment, err := r.db.QueryContext(ctx, queryComment, searchQuery, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, fmt.Errorf("Comment Search: %w", err)
	}

	var comments []comment.Comment
	var errFormat error
	for rowsComment.Next() {
		var comment comment.Comment
		if err := rowsComment.Scan(&comment.Body, &comment.DateCreated, &comment.ThreadID); err != nil {
			return nil, fmt.Errorf("Comment Search: %w", err)
		}

		comment.DateCreated, errFormat = config.FormatTime(comment.DateCreated, middleware.GetTimezone(ctx))
		if errFormat != nil {
			return nil, errFormat
		}

		comments = append(comments, comment)
	}

	if err = rowsComment.Err(); err != nil {
		return nil, fmt.Errorf("Comment Search: %w", err)
	}

	return comments, nil
}
