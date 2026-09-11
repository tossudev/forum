package thread

import (
	"database/sql"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) GetByCategory(id int) ([]Thread, error) {
	var threads []Thread

	rows, err := r.db.Query("SELECT * FROM threads WHERE category_id = ?", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (r *ThreadRepository) GetByID(id int) (*Thread, error) {
	query := "SELECT (id, title, body, date_created) FROM threads WHERE id = ?"
	var thread *Thread
	err := r.db.QueryRow(query, id).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated)
	if err != nil {
		return nil, err
	}

	return thread, nil

}

func (r *ThreadRepository) Create(threadRequest *Thread) (*Thread, error) {
	query := "INSERT INTO threads (title, body, date_created) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, threadRequest.Title, threadRequest.Body, threadRequest.DateCreated)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var newThread *Thread
	newThread, err = r.GetByID(int(id))
	if err != nil {
		return nil, err
	}
	return newThread, nil
}
