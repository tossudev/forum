package handlers

import (
	"database/sql"
	"forum/internal/entities/thread"
)

func initHandlers(db *sql.DB, validator *validator.Validate)  *ThreadHandler {
	threadRepo := thread.NewThreadRepo(db)
	threadService := thread.NewThreadService(threadRepo)
	threadHandler := thread.NewThreadHandler(threadService, validator)

	return threadHandler
}