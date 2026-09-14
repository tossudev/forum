package comment

import (
	"context"
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, req *Comment) (Comment, error) {

	var comment Comment

	query := `INSERT INTO comments body, date_created VALUES(?, ?) RETURNING id, body, date_created`
	err := r.db.QueryRowContext(ctx, query, req.Body, req.DateCreated).Scan(&comment.ID, &comment.Body, &comment.DateCreated)
	if err != nil {
		//handle error
	}

	return comment, nil
}

func (r *CommentRepository) GetByID(ctx context.Context, id int) (Comment, error) {

	//query := `SELECT (id, body, date_created) FROM comments WHERE id = ? RETURNING (id, body, date_created) `

	return Comment{}, nil

}

func (r *CommentRepository) GetByThread(id int) error {
	return nil
}
