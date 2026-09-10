package thread

import (
	"database/sql"
)

type ThreadRepository struct {
	db *sql.DB
}

func ThreadRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

func (r *ThreadRepository) GetThreadsByCategory() {
	
}