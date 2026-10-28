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

	//check if requested page is larger than last possible page based on number of threads & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("ThreadGetByCategory: %w", errs.ErrNotFound)
	}

	var threads []Thread

	query := `SELECT id, title, body, date_created, author_id, category_id
	FROM threads WHERE category_id = ? ORDER BY date_created ASC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, id, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("Thread GetByCategory: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
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

	//check if requested page is larger than last possible page based on number of comments & page size
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, nil, totalCount, fmt.Errorf("DisplayThreadPage: %w", pagination.ErrPageOverTotal) // totalCount needed for calculating the last page
	}

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

func (r *ThreadRepository) getThreadDisplay(ctx context.Context, threadID int, username string) (*ThreadDisplay, error) {
	query := `SELECT t.id, title, body, t.date_created, t.author_id, t.category_id, u.username, u.avatar,
				SUM(CASE WHEN tl.thread_like = 1 THEN 1 ELSE 0 END)	AS num_likes,
				SUM(CASE WHEN tl.thread_like = 0 THEN 1 ELSE 0 END)	AS num_dislikes,
				MAX(CASE WHEN tl.thread_like = 1 AND lu.username = ? THEN 1 ELSE 0 END)	    AS user_liked,
				MAX(CASE WHEN tl.thread_like = 0 AND lu.username = ? THEN 1 ELSE 0 END) 	AS user_disliked
			FROM threads AS t
			INNER JOIN users AS u 						ON u.id = t.author_id
			LEFT JOIN thread_likes AS tl 				ON tl.thread_id = t.id
			LEFT JOIN users AS lu 						ON lu.id = tl.user_id
			WHERE t.id = ? 
			GROUP BY t.id`

	var threadDisplay ThreadDisplay
	//TODO: Modify query to get author avatar
	row := r.db.QueryRowContext(ctx, query, username, username, threadID)
	err := row.Scan(&threadDisplay.Thread.ID,
		&threadDisplay.Thread.Title,
		&threadDisplay.Thread.Body,
		&threadDisplay.Thread.DateCreated,
		&threadDisplay.Thread.AuthorID,
		&threadDisplay.Thread.CategoryID,
		&threadDisplay.AuthorName,
		&threadDisplay.AuthorAvatar,
		&threadDisplay.LikeData.NumLikes,
		&threadDisplay.LikeData.NumDislikes,
		&threadDisplay.LikeData.LoggedUserLiked,
		&threadDisplay.LikeData.LoggedUserDisliked)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("getThreadDisplay: %w", err)
	}
	threadDisplay.Thread.DateCreated, err = timeformat.FormatTime(threadDisplay.Thread.DateCreated, timeformat.GetTimezone(ctx))
	if err != nil {
		return nil, fmt.Errorf("getThreadDisplay: %w", err)
	}
	return &threadDisplay, nil
}

func (r *ThreadRepository) getCommentsDisplay(ctx context.Context, threadID int, username string, paginationData pagination.Pagination) ([]CommentDisplay, error) {
	commentsQuery := `SELECT c.id, c.body, c.thread_id, u.username, u.avatar, c.date_created,
			SUM(CASE WHEN cl.comment_like = 1 THEN 1 ELSE 0 END)	AS num_likes,
			SUM(CASE WHEN cl.comment_like = 0 THEN 1 ELSE 0 END)	AS num_dislikes,
			MAX(CASE WHEN cl.comment_like = 1 AND lu.username = ? THEN 1 ELSE 0 END)	AS user_liked,
			MAX(CASE WHEN cl.comment_like = 0 AND lu.username = ? THEN 1 ELSE 0 END) 	AS user_disliked
		FROM comments AS c
		INNER JOIN users AS u			ON u.id = c.author_id
		LEFT JOIN comment_likes AS cl 	ON cl.comment_id = c.id
		LEFT JOIN users AS lu 			ON lu.id = cl.user_id
		WHERE c.thread_id = ? 
		GROUP BY c.id
		ORDER BY c.date_created ASC
		LIMIT ? OFFSET ?;`
	// values returned: id, body, threadID, username (author name), date created, number of likes/number of dislikes, which comments did the user like/dislike
	// scan said values into a slice and return
	//TODO: Get author avatar

	commentRows, err := r.db.QueryContext(ctx, commentsQuery, username, username, threadID, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, fmt.Errorf("getCommentsDisplay: %w", err)
	}
	defer commentRows.Close()
	var commentDisplays []CommentDisplay

	// fill the values to a commentdisplay and append to the slice above
	for commentRows.Next() {
		var cd CommentDisplay
		err := commentRows.Scan(&cd.Comment.ID,
			&cd.Comment.Body,
			&cd.Comment.ThreadID,
			&cd.AuthorName,
			&cd.AuthorAvatar,
			&cd.Comment.DateCreated,
			&cd.LikeData.NumLikes,
			&cd.LikeData.NumDislikes,
			&cd.LikeData.LoggedUserLiked,
			&cd.LikeData.LoggedUserDisliked)
		if err != nil {
			return nil, fmt.Errorf("getCommentsDisplay: %w", err)
		}

		cd.Comment.DateCreated, err = timeformat.FormatTime(cd.Comment.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, fmt.Errorf("getCommentsDisplay: %w", err)
		}
		commentDisplays = append(commentDisplays, cd)
	}

	if commentRows.Err() != nil {
		return nil, fmt.Errorf("getCommentsDisplay: %w", err)
	}
	return commentDisplays, nil
}
