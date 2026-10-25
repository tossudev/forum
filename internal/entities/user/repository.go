package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"forum/internal/database"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/timeformat"
)

type UserRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) AddUser(ctx context.Context, user *User) (*User, error) {
	query := `INSERT INTO users (username, email, password_hash, password_salt, date_created, role_id)
 VALUES (?, ?, ?, ?, ?, ?);`

	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password.Hash, user.Password.Salt, time.Now().Unix(), user.RoleID)
	if err != nil {
		if database.IsUniqueErr(err) {
			return nil, fmt.Errorf("%w: username or email already exists", errs.ErrDuplicate)
		}
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("adding user: getting last insert ID: %w", err)
	}

	return r.GetUserByID(ctx, int(userID))
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*User, error) {
	query := `SELECT id, username, email, password_hash, password_salt, date_created, role_id
	FROM users WHERE id = ?;`

	var user User
	var createdAt int
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password.Hash, &user.Password.Salt, &createdAt, &user.RoleID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: user not found", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserByID: scanning row: %w", err)
	}

	user.CreatedAt = time.Unix(int64(createdAt), 0)

	return &user, nil

}

func (r *UserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	// Query returns a single row containing 1 (true) or 0 (false)
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?);`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		return exists, fmt.Errorf("checking username exists: %w", err)
	}

	return exists, nil
}

func (r *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	// Query returns a single row containing 1 (true) or 0 (false)
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?);`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return exists, fmt.Errorf("checking email exists: %w", err)
	}

	return exists, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, username, email, password_hash, password_salt, date_created, role_id
	FROM users WHERE email = ?;`

	var user User
	var createdAt int
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password.Hash, &user.Password.Salt, &createdAt, &user.RoleID); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("%w: invalid credentials", errs.ErrUnauthorized)
		}
		return User{}, fmt.Errorf("GetUserByEmail: scanning row: %w", err)
	}

	user.CreatedAt = time.Unix(int64(createdAt), 0)

	return user, nil
}

// GetThreadsByUserID gets all the threads posted by a user within the given pagination parameters.
// Returns the threads, the total number of threads found, and a potential error
func (r *UserRepo) GetThreadsByUserID(ctx context.Context, userID int, paginationData pagination.Pagination) ([]ProfileThread, int, error) {

	//get count of threads created by user
	queryCount := `SELECT COUNT(*) FROM threads WHERE author_id = ?;`
	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	if totalCount == 0 {
		return nil, 0, nil
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", errs.ErrNotFound)
	}

	//get threads created by user on requested page
	var threads []ProfileThread

	query := `SELECT id, title, body, date_created, author_id, category_id FROM threads 
	WHERE author_id = ? ORDER BY date_created DESC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, userID, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
		}

		thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
		}

		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	return threads, totalCount, nil
}

// GetLikedThreadsByUserID gets all the threads liked by a user within the given pagination parameters.
// Returns the threads, the total number of threads found, and a potential error.
func (r *UserRepo) GetLikedThreadsByUserID(ctx context.Context, userID int, paginationData pagination.Pagination) ([]ProfileThread, int, error) {

	//get count of threads created by user
	queryCount := `SELECT COUNT(*) FROM thread_likes WHERE thread_like = true AND user_id = ?;`
	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetLikedThreadsByUserID: %w", err)
	}

	if totalCount == 0 {
		return nil, 0, nil
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetLikedThreadsByUserID: %w", errs.ErrNotFound)
	}

	var threads []ProfileThread

	query := `SELECT id, title, body, date_created, author_id, category_id
	FROM threads AS t
	INNER JOIN thread_likes AS tl ON t.id = tl.thread_id
	WHERE tl.thread_like = true AND tl.user_ID = ?
	ORDER BY date_created DESC
	LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, userID, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
		}
		thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
		}

		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	return threads, totalCount, nil
}

func (r *UserRepo) GetCommentsByUserID(ctx context.Context, userID int, pagination pagination.Pagination) ([]ProfileComment, int, int, error) {

	//check if page number is beyond possible pages
	//if so, make page number equal last page
	page, err := r.checkPage(ctx, userID, pagination)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetCommentsByUserID: %w", err)
	}
	pagination.Page = page

	var comments []ProfileComment
	var totalCount int

	query := `SELECT COUNT(*) OVER() AS total_count, id, body, date_created, thread_id, author_id
	FROM comments
	WHERE author_id = ?
	ORDER BY date_created DESC 
	LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, userID, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetCommentsByUserID: %w", err)
	}

	for rows.Next() {
		var comment ProfileComment
		if err := rows.Scan(&totalCount, &comment.ID, &comment.Body, &comment.DateCreated, &comment.ThreadID, &comment.AuthorID); err != nil {
			return nil, 0, 0, fmt.Errorf("GetCommentsByUserID: %w", err)
		}

		comment.DateCreated, err = timeformat.FormatTime(comment.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, 0, fmt.Errorf("GetCommentsByUserID: %w", err)
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("GetCommentsByUserID: %w", err)
	}

	return comments, totalCount, pagination.Page, nil
}

func (r *UserRepo) checkPage(ctx context.Context, userID int, paginationData pagination.Pagination) (int, error) {
	query := `SELECT COUNT(*) OVER() AS total_count
	FROM comments
	WHERE author_id = ?;`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, nil
	}

	if paginationData.Page > pagination.CalcTotalPages(count, paginationData.PageSize) {
		return pagination.CalcTotalPages(count, paginationData.PageSize), nil
	}

	return paginationData.Page, nil
}
