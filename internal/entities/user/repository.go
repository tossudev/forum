package user

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/database"
	"forum/internal/errs"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) AddUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password_hash, password_salt, date_created, role_id)
 VALUES (?, ?, ?, ?, ?, ?);`

	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password.Hash, user.Password.Salt, time.Now().Unix(), user.RoleID)
	if err != nil {
		if database.IsUniqueErr(err) {
			return fmt.Errorf("%w: username or email already exists", errs.ErrDuplicate)
		}
		return err
	}

	return nil
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
	query := `SELECT username, email, password_hash, password_salt, date_created, role_id
	FROM users WHERE email = ?;`

	var user User
	var createdAt int
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Username, &user.Email, &user.Password.Hash, &user.Password.Salt, &createdAt, &user.RoleID); err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("%w: invalid credentials", errs.ErrsUnauthorized)
		}
		return User{}, fmt.Errorf("GetUserByEmail: scanning row: %w", err)
	}

	user.CreatedAt = time.Unix(int64(createdAt), 0)

	return user, nil
}
