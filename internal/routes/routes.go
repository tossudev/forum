package routes

import (
	"net/http"
)

func GetRoutes(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /users/register", app.UserHandler.RegisterUser)
	mux.HandleFunc("POST /users/login", app.UserHandler.Login)

	// Category routes
	mux.HandleFunc("POST /categories", app.CategoryHandler.Create)

	// Thread routes
	mux.HandleFunc("GET /categories/{id}", app.ThreadHandler.GetByCategory)
	mux.HandleFunc("POST /threads", app.ThreadHandler.Create)
	mux.HandleFunc("GET /threads/{id}", app.ThreadHandler.GetByID)

	//Comment routes
	mux.HandleFunc("GET /comments/{id}", app.CommentHandler.GetByID)
	mux.HandleFunc("GET /threads/{id}/comments", app.CommentHandler.GetByThread)
	mux.HandleFunc("POST /threads/{id}/comments", app.CommentHandler.Create)
	//TODO: on comments --> force delete

	// Like routes
	mux.HandleFunc("POST /like/thread", app.LikeHandler.LikeThread)
	mux.HandleFunc("POST /like/comment", app.LikeHandler.LikeComment)

	return mux

}
