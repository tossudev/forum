package comment

import (
	"context"
	"database/sql"
	"errors"

	"forum/internal/errs"
	"forum/internal/pagination"
)

type CommentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, req *Comment) (Comment, error) {

	var comment Comment

	query := `INSERT INTO comments (body, date_created, thread_id, author_id) 
			  VALUES(?, ?, ?, ?) 
			  RETURNING id, body, date_created, thread_id, author_id`

	row := r.db.QueryRowContext(ctx, query, req.Body, req.DateCreated, req.ThreadID, req.AuthorID)
	err := row.Scan(&comment.ID, &comment.Body, &comment.DateCreated, &comment.ThreadID, &comment.AuthorID)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (r *CommentRepository) GetByID(ctx context.Context, id int) (Comment, error) {

	comment := Comment{}
	query := `SELECT id, body, date_created, thread_id, author_id FROM comments WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.Body, &comment.DateCreated, &comment.ThreadID, &comment.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Comment{}, errs.ErrNotFound
		}
		return Comment{}, err
	}
	return comment, nil
}

func (r *CommentRepository) GetByThread(ctx context.Context, id int, pag pagination.Pagination) ([]Comment, error) {

  query := `SELECT id,
 				   body,
				   date_created,
				   thread_id,
				   author_id 
				   FROM comments WHERE thread_id = ?
				   ORDER BY date_created
				   ASC LIMIT ? OFFSET ?;`

	comments := []Comment{}

	rows, err := r.db.QueryContext(ctx, query, id, pag.Limit(), pag.Offset())
	if err != nil { return nil, err }
	defer rows.Close()

	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.ID, &comment.Body, &comment.DateCreated, &comment.ThreadID, &comment.AuthorID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil { 
		return nil, err 
	}
	
	return comments, nil 
}

func (r *CommentRepository) Delete(ctx context.Context, id int) error {

	query := `DELETE FROM comments WHERE id = ?` 
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil { return err }
	
	amount, err := res.RowsAffected()
	if err != nil { return err }
	if amount == 0 { return errs.ErrNotFound }

	return nil 
}
