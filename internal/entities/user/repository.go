package user

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/errs"
	"forum/internal/pagination"
	"time"
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
// Returns the threads, the total number of threads found, and a potential error.
func (r *UserRepo) GetThreadsByUserID(ctx context.Context, userID int, pagination pagination.Pagination) ([]ProfileThread, int, error) {
	var threads []ProfileThread
	var totalCount int

	query := `SELECT COUNT(*) OVER() AS total_count, id, title, body, date_created, author_id, category_id FROM threads 
	WHERE author_id = ? ORDER BY date_created DESC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, userID, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&totalCount, &thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
		}

		thread.DateCreated, err = config.FormatTime(thread.DateCreated)
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
func (r *UserRepo) GetLikedThreadsByUserID(ctx context.Context, userID int, pagination pagination.Pagination) ([]ProfileThread, int, error) {
	var threads []ProfileThread
	var totalCount int

	query := `SELECT COUNT(*) OVER() AS total_count, id, title, body, date_created, author_id, category_id FROM threads 
	INNER JOIN thread_likes AS likes ON likes.user_id = ?
	WHERE likes.thread_like = true ORDER BY date_created DESC LIMIT ? OFFSET ?;`

	rows, err := r.db.QueryContext(ctx, query, userID, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
	}

	for rows.Next() {
		var thread ProfileThread
		if err := rows.Scan(&totalCount, &thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, 0, fmt.Errorf("GetThreadsByUserID: %w", err)
		}

		thread.DateCreated, err = config.FormatTime(thread.DateCreated)
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
