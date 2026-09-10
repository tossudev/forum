package handlers

import (
	"forum/internal/database"
	"forum/internal/entities/thread"
)

type App struct {
	ThreadHandler *thread.ThreadHandler
}

func initApp(threadHandler *ThreadHandler) *App {

	app := App{
		ThreadHandler: threadHandler
	}

	return &app
}