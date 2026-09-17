package category

import (
	"database/sql"
	"context"
	"fmt"

)

type CategoryRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*Category, error) {
	query := "SELECT id, name FROM categories WHERE id = ?;"
	var category Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		fmt.Println("here", id)
		return nil, err
	}
	
	return &category, nil
}
 
func (r *CategoryRepository) Create(ctx context.Context, category *Category) (*Category, error) {
	query := "INSERT INTO categories (name) VALUES (?);"	
	result, err := r.db.ExecContext(ctx, query, category.Name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var newCategory *Category
	newCategory, err = r.GetByID(ctx, int(id))
	if err != nil {
		return nil, err
	}
	return newCategory, nil
}