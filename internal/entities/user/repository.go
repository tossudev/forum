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

	createdAt := time.Now().UTC().Format(timeformat.TimeFormat)
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password.Hash, user.Password.Salt, createdAt, user.RoleID)
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password.Hash, &user.Password.Salt, &createdAt, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: user not found", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("GetUserByID: scanning row: %w", err)
	}

	user.CreatedAt, err = timeformat.FormatTime(user.CreatedAt, timeformat.GetTimezone(ctx))

	return &user, nil

}

func (r *UserRepo) GetUserProfileByUsername(ctx context.Context, username string) (*UserProfile, error) {
	query := `SELECT username, date_created
	FROM users WHERE username = ?;`

	var up UserProfile
	err := r.db.QueryRowContext(ctx, query, username).Scan(&up.Username, &up.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("GetUserProfileByUsername: %w", err)
	}

	up.CreatedAt, err = timeformat.FormatTime(up.CreatedAt, timeformat.GetTimezone(ctx))

	return &up, nil
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
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password.Hash, &user.Password.Salt, &user.CreatedAt, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("%w: invalid credentials", errs.ErrUnauthorized)
		}
		return User{}, fmt.Errorf("GetUserByEmail: scanning row: %w", err)
	}

	user.CreatedAt, err = timeformat.FormatTime(user.CreatedAt, timeformat.GetTimezone(ctx))

	return user, nil
}

// GetThreadsByUsername gets all the threads posted by a user within the given pagination parameters.
// Returns the threads, the total number of threads found, and a potential error
func (r *UserRepo) GetThreadsByUsername(ctx context.Context, username string, paginationData pagination.Pagination) ([]ProfileThread, int, error) {

	//get count of threads created by user
	queryCount := `SELECT COUNT(*) FROM threads
	INNER JOIN users ON users.id = threads.author_id WHERE users.username = ?;`

	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, username).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
	}

	if totalCount == 0 {
		return nil, 0, nil
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", errs.ErrNotFound)
	}

	//get threads created by user on requested page
	var threads []ProfileThread

	query := `SELECT t.id, t.title, t.body, t.date_created, t.author_id, t.category_id
	FROM threads AS t INNER JOIN users ON users.id = t.author_id
	WHERE users.username = ? ORDER BY t.date_created DESC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, username, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
		}

		thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
		}

		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
	}

	return threads, totalCount, nil
}

// GetLikedThreadsByUsername gets all the threads liked by a user within the given pagination parameters.
// Returns the threads, the total number of threads found, and a potential error.
func (r *UserRepo) GetLikedThreadsByUsername(ctx context.Context, username string, paginationData pagination.Pagination) ([]ProfileThread, int, error) {

	//get count of threads liked by user
	queryCount := `SELECT COUNT(*) FROM thread_likes
	INNER JOIN users ON users.id = thread_likes.user_id 
	WHERE thread_like = true AND users.Username = ?;`
	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, username).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetLikedThreadsByUsername: %w", err)
	}

	if totalCount == 0 {
		return nil, 0, nil
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetLikedThreadsByUsername: %w", errs.ErrNotFound)
	}

	var threads []ProfileThread

	query := `SELECT t.id, t.title, t.body, t.date_created, t.author_id, t.category_id
	FROM threads AS t
	INNER JOIN thread_likes AS tl ON t.id = tl.thread_id
	INNER JOIN users ON users.id = tl.user_id
	WHERE tl.thread_like = true AND users.username = ?
	ORDER BY t.date_created DESC
	LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, username, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
		}
		thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
		}

		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUsername: %w", err)
	}

	return threads, totalCount, nil
}

func (r *UserRepo) GetCommentsByUsername(ctx context.Context, username string, paginationData pagination.Pagination) ([]ProfileComment, int, error) {

	//get count of comments created by user
	queryCount := `SELECT COUNT(*) FROM comments 
	INNER JOIN users ON users.id = comments.author_id WHERE users.username = ?;`

	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, username).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetCommentsByUsername: %w", err)
	}

	if totalCount == 0 {
		return nil, 0, nil
	}

	//check if page number is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetCommentsByUsername: %w", errs.ErrNotFound)
	}

	var comments []ProfileComment

	query := `SELECT c.id, c.body, c.date_created, c.thread_id, c.author_id
	FROM comments AS c
	INNER JOIN users ON users.id = c.author_id
	WHERE users.username = ?
	ORDER BY c.date_created DESC 
	LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, username, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetCommentsByUsername: %w", err)
	}

	for rows.Next() {
		var comment ProfileComment
		if err := rows.Scan(&comment.ID, &comment.Body, &comment.DateCreated, &comment.ThreadID, &comment.AuthorID); err != nil {
			return nil, 0, fmt.Errorf("GetCommentsByUsername: %w", err)
		}

		comment.DateCreated, err = timeformat.FormatTime(comment.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("GetCommentsByUsername: %w", err)
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetCommentsByUsername: %w", err)
	}

	return comments, totalCount, nil
}

// Returns the threads, the total number of threads found, and a potential error.
func (r *UserRepo) GetDislikedThreadsByUsername(ctx context.Context, username string, paginationData pagination.Pagination) ([]ProfileThread, int, error) {

	//get count of threads disliked by user
	queryCount := `SELECT COUNT(*) FROM thread_likes
	INNER JOIN users ON users.id = thread_likes.user_id 
	WHERE thread_like = false AND users.username = ?;`

	var totalCount int
	err := r.db.QueryRowContext(ctx, queryCount, username).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetDisikedThreadsByUsername: %w", err)
	}

	if totalCount == 0 {
		return nil, 0, nil
	}

	//check if requested page is larger than last possible page based on number of results & page size
	//if so, return err
	err = pagination.CheckPageTooLarge(paginationData, totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("GetDisikedThreadsByUsername: %w", errs.ErrNotFound)
	}

	var threads []ProfileThread

	query := `SELECT t.id, t.title, t.body, t.date_created, t.author_id, t.category_id
	FROM threads AS t
	INNER JOIN thread_likes AS tl ON t.id = tl.thread_id
	INNER JOIN users ON users.id = tl.user_id
	WHERE tl.thread_like = false AND users.username = ?
	ORDER BY t.date_created DESC
	LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, username, paginationData.Limit(), paginationData.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetDisikedThreadsByUsername: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetDisikedThreadsByUsername: %w", err)
		}
		thread.DateCreated, err = timeformat.FormatTime(thread.DateCreated, timeformat.GetTimezone(ctx))
		if err != nil {
			return nil, 0, fmt.Errorf("GetDisikedThreadsByUsername: %w", err)
		}

		threads = append(threads, thread)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetDisikedThreadsByUsername: %w", err)
	}

	return threads, totalCount, nil
}
