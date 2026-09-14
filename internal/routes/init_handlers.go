package routes

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/middleware"

	//"forum/internal/middleware"
	"github.com/go-playground/validator/v10"
)

type App struct {
	UserHandler   *user.UserHandler
	ThreadHandler *thread.ThreadHandler
}

func InitHandlers(db *sql.DB, validate *validator.Validate) http.Handler {
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)
	threadHandler := thread.NewHandler(threadService, validate)

	app := App{
		UserHandler:   userHandler,
		ThreadHandler: threadHandler,
	}

	mux := GetRoutes(&app)

	// Middleware
	handler := middleware.Timeout(5 * time.Second)(mux)
	handler = middleware.Logger(handler)
	handler = middleware.RecoverPanic(handler)

	return handler
}
