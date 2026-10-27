package routes

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/entities/category"
	"forum/internal/entities/comment"
	"forum/internal/entities/image"
	"forum/internal/entities/like"
	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/middleware"
	"forum/internal/render"
	"forum/internal/search"
	"forum/internal/session"
	"forum/internal/timeformat"
)

type App struct {
	UserHandler     *user.UserHandler
	CategoryHandler *category.CategoryHandler
	ThreadHandler   *thread.ThreadHandler
	CommentHandler  *comment.CommentHandler
	LikeHandler     *like.LikeHandler
	ImageHandler    *image.ImageHandler
	SearchHandler   *search.SearchHandler
}

func InitHandlers(db *sql.DB, renderer *render.Renderer) http.Handler {
	sessionRepo := session.NewRepo(db)
	sm := session.NewSessionManager(sessionRepo, "session_token", 30*24*time.Hour) // session expires after a 30 days of inactivity

	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, sm, renderer)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService, renderer)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)
	threadHandler := thread.NewHandler(threadService, renderer)

	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, renderer, threadService)

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)
	likeHandler := like.NewHandler(likeService)

	imageRepo := image.NewRepository(db)
	imageService := image.NewService(imageRepo)
	imageHandler := image.NewHandler(imageService)

	searchRepo := search.NewRepository(db)
	searchService := search.NewService(searchRepo)
	searchHandler := search.NewHandler(searchService, renderer)

	app := App{
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		ThreadHandler:   threadHandler,
		CommentHandler:  commentHandler,
		LikeHandler:     likeHandler,
		ImageHandler:    imageHandler,
		SearchHandler:   searchHandler,
	}

	mux := GetRoutes(&app, sm)

	// Middleware
	handler := middleware.Timeout(5 * time.Second)(mux)
	handler = userHandler.GetUserMiddleware(handler)
	handler = sm.Authenticate(handler)
	handler = timeformat.Timezone(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RecoverPanic(handler)

	return handler
}
