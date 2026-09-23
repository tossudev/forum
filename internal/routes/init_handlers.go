package routes

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/entities/category"
	"forum/internal/entities/comment"
	"forum/internal/entities/like"
	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/middleware"
	"forum/internal/render"
	"forum/internal/session"

	//"forum/internal/middleware"
	"github.com/go-playground/validator/v10"
)

type App struct {
	UserHandler     *user.UserHandler
	CategoryHandler *category.CategoryHandler
	ThreadHandler   *thread.ThreadHandler
	CommentHandler  *comment.CommentHandler
	LikeHandler     *like.LikeHandler
}

func InitHandlers(db *sql.DB, validate *validator.Validate, renderer *render.Renderer) http.Handler {
	sessionRepo := session.NewRepo(db)
	sm := session.NewSessionManager(sessionRepo, "session_token", 30*24*time.Hour) // session expires after a 30 days of inactivity

	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, sm, renderer)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService, validate, renderer)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)
	threadHandler := thread.NewHandler(threadService, validate, renderer)

	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, validate)

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)
	likeHandler := like.NewHandler(likeService, validate)

	app := App{
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		ThreadHandler:   threadHandler,
		CommentHandler:  commentHandler,
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
