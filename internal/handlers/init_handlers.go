package handlers

import (
	"database/sql"
	"forum/internal/entities/thread"
)

func initHandlers(db *sql.DB, validator *validator.Validate)  *ThreadHandler {
	threadRepo := thread.NewRepo(db)
	threadService := thread.NewService(threadRepo)
	threadHandler := thread.NewHandler(threadService, validator)

	return threadHandler
}