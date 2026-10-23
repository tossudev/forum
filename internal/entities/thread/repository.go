package thread

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

// Returns threads, the total number of threads, and a potential error.
func (r *ThreadRepository) GetByCategory(ctx context.Context, id int, paginationData pagination.Pagination) ([]Thread, int, int, error) {

	//get count of threads in requested category
	queryCount := `SELECT COUNT (*) FROM threads WHERE category_id = ?;`
	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, id).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, set page to last possible page
	page := pagination.CheckPageTooLarge(paginationData, totalCount)
	paginationData.Page = page

	var threads []Thread

	query := `SELECT COUNT(*) OVER() AS total_count, id, title, body, date_created, author_id, category_id
	FROM threads WHERE category_id = ? ORDER BY date_created ASC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, id, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&totalCount, &thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, 0, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		queryCount := "SELECT COUNT(*) FROM comments WHERE thread_id = ?"
		if err := r.db.QueryRowContext(ctx, queryCount, thread.ID).Scan(&thread.ReplyCount); err != nil {
			return nil, 0, 0, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	return threads, totalCount, page, nil
}

func (r *ThreadRepository) GetByID(ctx context.Context, id int) (*Thread, error) {
	query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE id = ?;"
	var thread Thread
	err := r.db.QueryRowContext(ctx, query, id).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("Thread GetByID: %w", err)
	}

	return &thread, nil

}

func (r *ThreadRepository) GetCategoryName(ctx context.Context, categoryID int) (string, error) {
	query := "SELECT name FROM categories WHERE id = ?;"

	var categoryName string
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&categoryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("GetCategoryName: %w", errs.ErrNotFound)
		}
		return "", fmt.Errorf("GetCategoryName: %w", err)
	}

	return categoryName, nil
}

func (r *ThreadRepository) Create(ctx context.Context, thread *Thread) (int, error) {
	query := "INSERT INTO threads (title, body, date_created, author_id, category_id) VALUES (?, ?, ?, ?, ?);"
	result, err := r.db.ExecContext(ctx, query, thread.Title, thread.Body, thread.DateCreated, thread.AuthorID, thread.CategoryID)
	if err != nil {
		return 0, fmt.Errorf("Thread Create: %w", err)
	}

	threadID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Thread Create: getting last insert id: %w", err)
	}

	return int(threadID), nil
}

func (r *ThreadRepository) Delete(ctx context.Context, threadID int, userID int) error {
	thread, err := r.GetByID(ctx, threadID)
	if err != nil {
		return fmt.Errorf("Thread search: %w", err)
	}

	if userID != thread.AuthorID {
		return fmt.Errorf("%w: user is not thread creator", errs.ErrUnauthorized)
	}

	query := "DELETE from threads WHERE id = ?;"
	_, err = r.db.ExecContext(ctx, query, threadID)
	if err != nil {
		return fmt.Errorf("Thread delete: %w", err)
	}

	return nil
}

// DisplayThreadPage returns *ThreadData, containing all the data needed for displaying a single thread page (including comments and associated likes).
func (r *ThreadRepository) DisplayThreadPage(ctx context.Context, threadID int, username string) (*ThreadData, error) {
	// Get thread data along with thread author name
	threadQuery := `SELECT threads.id, title, body, threads.date_created, threads.author_id, threads.category_id, username
	FROM users INNER JOIN threads ON users.id = threads.author_id WHERE threads.id = ?;`
	//query := `SELECT threads.id, title, body, date_created, username, thread_like, thread_likes.user_id FROM threads INNER JOIN users ON users.id = threads.author_id LEFT JOIN thread_likes ON threads.id = thread_likes.thread_id;`
	//query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE id = ?;"

	var authorName string
	var thread Thread
	err := r.db.QueryRowContext(ctx, threadQuery, threadID).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID, &authorName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	thread.DateCreated, err = config.FormatTime(thread.DateCreated)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	// Get thread likes data
	//this is a separate query because otherwise there would be a distinct row in the above query for every like on every thread
	threadLikeQuery := `SELECT username, thread_like FROM thread_likes INNER JOIN users ON users.id = thread_likes.user_id WHERE thread_id = ?`
	threadLikeData, err := r.getLikeData(ctx, threadLikeQuery, threadID, username)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	//get all comments for this thread plus username of comment author
	threadCommentsQuery := `SELECT comments.id, comments.body, comments.thread_id, username, comments.date_created FROM comments INNER JOIN users ON comments.author_id = users.id WHERE comments.thread_id = ?;`
	commentRows, err := r.db.QueryContext(ctx, threadCommentsQuery, threadID)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	var comments []CommentData

	// Get the data for each comment in the thread
	for commentRows.Next() {
		var comment CommentData
		if err := commentRows.Scan(&comment.ID, &comment.Body, &comment.ThreadID, &comment.AuthorName, &comment.DateCreated); err != nil {
			return nil, fmt.Errorf("DisplayThreadPage: %w", err)
		}

		comment.DateCreated, err = config.FormatTime(comment.DateCreated)
		if err != nil {
			return nil, fmt.Errorf("DisplayThreadPage: %w", err)
		}

		// Get comment likes data
		commentLikeQuery := `SELECT username, comment_like FROM comment_likes INNER JOIN users ON users.id = comment_likes.user_id WHERE comment_id = ?;`
		commentLikeData, err := r.getLikeData(ctx, commentLikeQuery, comment.ID, username)
		if err != nil {
			return nil, fmt.Errorf("DisplayThreadPage: %w", err)
		}
		comment.LikeData = *commentLikeData
		comments = append(comments, comment)
	}

	data := ThreadData{
		Thread:     thread,
		AuthorName: authorName,
		LikeData:   *threadLikeData,
		Comments:   comments,
	}

	return &data, nil
}

// getLikeData gets the LikeData for a single thread or comment, where id can be either the thread ID or comment ID.
// The query provided should return rows with 1) the username of the like/dislike, and 2) the like status (bool), in that order.
func (r *ThreadRepository) getLikeData(ctx context.Context, query string, id int, username string) (*LikeData, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("getLikeData: %w", err)
	}

	var likeCount, dislikeCount int
	var loggedUserLike, loggedUserDislike bool
	for rows.Next() {
		var likeUser string
		var like bool
		if err := rows.Scan(&likeUser, &like); err != nil {
			return nil, fmt.Errorf("getLikeData: %w", err)
		}

		if like {
			likeCount++
			if likeUser == username {
				loggedUserLike = true
			}
		} else {
			dislikeCount++
			if likeUser == username {
				loggedUserDislike = true
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getLikeData: %w", err)
	}

	likeData := LikeData{
		NumLikes:           likeCount,
		NumDislikes:        dislikeCount,
		LoggedUserLiked:    loggedUserLike,
		LoggedUserDisliked: loggedUserDislike,
	}

	return &likeData, nil
}
