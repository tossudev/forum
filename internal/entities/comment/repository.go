package comment

import (
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(req Comment) (Comment, error) {
	return Comment{}, nil
}

func (r *CommentRepository) GetByThread(id int) error {
	return nil
}
