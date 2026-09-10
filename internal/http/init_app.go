package http

import (
	"forum/internal/database"
	"forum/internal/entities/thread"
)

type App struct {
	ThreadService *service.ThreadService
}

func initApp(db *sql.DB, validate *validator.Validate) *handlers.App {
	repo := &respository.Repo{DB: db}

	app := App{
		ThreadService: thread.NewThreadService(repo, validate)
	}

	return &app
}