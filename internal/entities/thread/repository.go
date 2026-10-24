package thread

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/timeformat"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

// Returns threads, the total number of threads, and a potential error.
func (r *ThreadRepository) GetByCategory(ctx context.Context, id int, paginationData pagination.Pagination) ([]Thread, int, error) {

	//get count of threads in requested category
	queryCount := `SELECT COUNT (*) FROM threads WHERE category_id = ?;`
	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, id).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("ThreadGetByCategory: %w", errs.ErrInvalidUserInput)
	}

	var threads []Thread

	query := `SELECT COUNT(*) OVER() AS total_count, id, title, body, date_created, author_id, category_id
	FROM threads WHERE category_id = ? ORDER BY date_created ASC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, id, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&totalCount, &thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		queryCount := "SELECT COUNT(*) FROM comments WHERE thread_id = ?"
		if err := r.db.QueryRowContext(ctx, queryCount, thread.ID).Scan(&thread.ReplyCount); err != nil {
			return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	return threads, totalCount, nil
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

// DisplayThreadPage gets the data needed for displaying a single thread page (including comments and associated likes).
// Also returns an int of the total comments count for pagination.
func (r *ThreadRepository) DisplayThreadPage(ctx context.Context, threadID int, username string, paginationData pagination.Pagination) (*ThreadDisplay, []CommentDisplay, int, error) {
	// Get count of comments in requested thread
	countQuery := `SELECT COUNT (*) FROM comments WHERE thread_id = ?;`
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, threadID).Scan(&totalCount)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	// TODO: Use pagination.CheckPageTooLarge once it has been updated to return error

	// Get thread data along with thread author name
	threadDisplay, err := r.getThreadDisplay(ctx, threadID, username)
	if err != nil {
		return nil, nil, 0, err
	}

	// Get comment data along with comment author data
	commentsDisplay, err := r.getCommentsDisplay(ctx, threadID, username, paginationData)
	if err != nil {
		return nil, nil, 0, err
	}

	return threadDisplay, commentsDisplay, totalCount, nil
}

// getThreadDisplay is a helper for getting all data for displaying a single thread post.
func (r *ThreadRepository) getThreadDisplay(ctx context.Context, threadID int, username string) (*ThreadDisplay, error) {
	// Get thread data along with thread author name
	query := `SELECT threads.id, title, body, threads.date_created, threads.author_id, threads.category_id, username
	FROM users INNER JOIN threads ON users.id = threads.author_id WHERE threads.id = ?;`

	//TODO: Modify query to get author avatar

	var authorName string
	var thread Thread
	err := r.db.QueryRowContext(ctx, query, threadID).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID, &authorName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("getThreadDisplay: %w", err)
	}

	thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
	if err != nil {
		return nil, fmt.Errorf("getThreadDisplay: %w", err)
	}

	// Get thread likes data
	// TODO: Reevaluate whether the queries can be optimized further instead of this being separate
	likeQuery := `SELECT username, thread_like FROM thread_likes INNER JOIN users ON users.id = thread_likes.user_id WHERE thread_id = ?`
	likeData, err := r.getLikeData(ctx, likeQuery, threadID, username)
	if err != nil {
		return nil, fmt.Errorf("getThreadDisplay: %w", err)
	}

	threadDisplay := ThreadDisplay{
		Thread:     thread,
		AuthorName: authorName,
		//AuthorAvatar: authorAvatar,
		LikeData: *likeData,
		//DisplayError: ?,
	}

	return &threadDisplay, nil
}

// getCommentsDisplay is a helper for getting all the comment data (paginated) for displaying comments.
func (r *ThreadRepository) getCommentsDisplay(ctx context.Context, threadID int, username string, paginationData pagination.Pagination) ([]CommentDisplay, error) {
	// Get all comments for this thread plus username of comment author
	commentsQuery := `SELECT c.id, c.body, c.thread_id, u.username, c.date_created
	FROM comments AS c INNER JOIN users AS u ON c.author_id = u.id WHERE c.thread_id = ?
	ORDER BY c.date_created ASC LIMIT ? OFFSET ?;`
	commentRows, err := r.db.QueryContext(ctx, commentsQuery, threadID, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, fmt.Errorf("getCommentsDisplay: %w", err)
	}

	var comments []CommentDisplay

	// Get the data for each comment in the thread
	for commentRows.Next() {
		var comment Comment
		var authorName string
		//TODO: Get author avator
		if err := commentRows.Scan(&comment.ID, &comment.Body, &comment.ThreadID, &authorName, &comment.DateCreated); err != nil {
			return nil, fmt.Errorf("getCommentsDisplay: %w", err)
		}

		comment.DateCreated, err = timeformat.FormatTime(comment.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, fmt.Errorf("getCommentsDisplay: %w", err)
		}

		// Get comment likes data
		// TODO: Refactor to combine this and above query
		likeQuery := `SELECT username, comment_like FROM comment_likes INNER JOIN users ON users.id = comment_likes.user_id WHERE comment_id = ?;`
		likeData, err := r.getLikeData(ctx, likeQuery, comment.ID, username)
		if err != nil {
			return nil, fmt.Errorf("getCommentsDisplay: %w", err)
		}

		commentDisplay := CommentDisplay{
			Comment:    comment,
			AuthorName: authorName,
			//AuthorAvatar: authorAvatar,
			LikeData: *likeData,
		}
		comments = append(comments, commentDisplay)
	}

	return comments, nil
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
