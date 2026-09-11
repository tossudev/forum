package routes

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/entities/thread"
	"forum/internal/middleware"

	//"forum/internal/middleware"
	"github.com/go-playground/validator/v10"
)

type App struct {
	ThreadHandler *thread.ThreadHandler
}

func InitHandlers(db *sql.DB, validate *validator.Validate) http.Handler {
	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)
	threadHandler := thread.NewHandler(threadService, validate)

	app := App{
		ThreadHandler: threadHandler,
	}

	mux := GetRoutes(&app)

	// Middleware
	handler := middleware.Timeout(5 * time.Second)(mux)
	handler = middleware.Logger(handler)
	handler = middleware.RecoverPanic(handler)

	return handler
}
