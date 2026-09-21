package thread

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/errs"
	"forum/internal/pagination"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) GetByCategory(ctx context.Context, id int, pagination pagination.Pagination) ([]Thread, error) {
	var threads []Thread

	query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE category_id = ? ORDER BY date_created ASC LIMIT ? OFFSET ?;"

	rows, err := r.db.QueryContext(ctx, query, id, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, fmt.Errorf("Thread GetByCategory: %w", err)
		}
		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	return threads, nil
}

func (r *ThreadRepository) GetByID(ctx context.Context, id int) (*Thread, error) {
	query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE id = ?;"
	var thread Thread
	err := r.db.QueryRowContext(ctx, query, id).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("Thread GetByID: %w", err)
	}

	return &thread, nil

}

func (r *ThreadRepository) Create(ctx context.Context, thread *Thread) error {
	query := "INSERT INTO threads (title, body, date_created, author_id, category_id) VALUES (?, ?, ?, ?, ?);"
	_, err := r.db.ExecContext(ctx, query, thread.Title, thread.Body, thread.DateCreated, thread.AuthorID, thread.CategoryID)
	if err != nil {
		return fmt.Errorf("Thread Create: %w", err)
	}

	return nil
}

func (r *ThreadRepository) Delete(ctx context.Context, threadID int, userID int) error {
	thread, err := r.GetByID(ctx, threadID)
	if err != nil {
		return fmt.Errorf("Thread search: %w", err)
	}

	threadCreator := thread.AuthorID
	if userID != threadCreator {
		return fmt.Errorf("%w: user is not thread creator", errs.ErrUnauthorized)
	}

	query := "DELETE from threads WHERE id = ?;"
	_, err = r.db.ExecContext(ctx, query, threadID)
	if err != nil {
		return fmt.Errorf("Thread delete: %w", err)
	}

	return nil
}
