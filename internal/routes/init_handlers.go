package routes

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/entities/category"
	"forum/internal/entities/like"
	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/middleware"
	"forum/internal/session"

	//"forum/internal/middleware"
	"github.com/go-playground/validator/v10"
)

type App struct {
	UserHandler     *user.UserHandler
	CategoryHandler *category.CategoryHandler
	ThreadHandler   *thread.ThreadHandler
	LikeHandler     *like.LikeHandler
}

func InitHandlers(db *sql.DB, validate *validator.Validate) http.Handler {
	sessionRepo := session.NewRepo(db)
	sm := session.NewSessionManager(sessionRepo, "session_token", 30*24*time.Hour) // session expires after a 30 days of inactivity

	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, sm)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService, validate)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)
	threadHandler := thread.NewHandler(threadService, validate)

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)
	likeHandler := like.NewHandler(likeService, validate)

	app := App{
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		ThreadHandler:   threadHandler,
		LikeHandler:     likeHandler,
	}

	mux := GetRoutes(&app, sm)

	// Middleware
	handler := middleware.Timeout(5 * time.Second)(mux)
	handler = sm.Authenticate(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RecoverPanic(handler)

	return handler
}
