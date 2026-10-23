package category

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/pagination"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll(ctx context.Context, pagination pagination.Pagination) ([]Category, error) {
	var categories []Category
	query := "SELECT id, name FROM categories ORDER BY id ASC LIMIT ? OFFSET ?;"
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, fmt.Errorf("Category GetAll: %w", err)
	}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("Category GetAll: %w", err)
		}

		queryCount := "SELECT COUNT(*) FROM threads WHERE category_id = ?"
		if err := r.db.QueryRowContext(ctx, queryCount, category.ID).Scan(&category.ThreadCount); err != nil {
			return nil, fmt.Errorf("Category GetAll: %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Category GetAll: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) Create(ctx context.Context, category *Category) error {
	query := "INSERT INTO categories (name) VALUES (?);"
	_, err := r.db.ExecContext(ctx, query, category.Name)
	if err != nil {
		return fmt.Errorf("Category Create: %w", err)
	}

	return nil
}
