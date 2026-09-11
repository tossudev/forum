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

func (r *ThreadRepository) GetByCategory() {
	
}