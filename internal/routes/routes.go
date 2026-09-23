package routes

import (
	"forum/internal/session"
	"net/http"
)

func GetRoutes(app *App, sm *session.SessionManager) *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /users/register", app.UserHandler.RegisterUser)
	mux.HandleFunc("POST /users/login", app.UserHandler.Login)
	mux.HandleFunc("POST /users/logout", app.UserHandler.Logout)

	// Category routes
	mux.HandleFunc("POST /categories", app.CategoryHandler.Create)
	mux.HandleFunc("GET /", app.CategoryHandler.GetAll)

	// Thread routes
	mux.HandleFunc("GET /categories/{id}", app.ThreadHandler.GetByCategory)
	mux.HandleFunc("POST /threads", app.ThreadHandler.Create)
	mux.HandleFunc("GET /threads/{id}", app.ThreadHandler.GetByID)
	mux.HandleFunc("DELETE /threads/{id}", app.ThreadHandler.Delete)

	//Comment routes
	mux.HandleFunc("GET /comments/{id}", app.CommentHandler.GetByID)
	mux.HandleFunc("GET /threads/{id}/comments", app.CommentHandler.GetByThread)
	mux.HandleFunc("POST /threads/{id}/comments", app.CommentHandler.Create)
	//TODO: fix the delete endpoint, test throughly
	//mux.HandleFunc("DELETE /comments/{id}", app.CommentHandler.Delete)

	// Like routes
	mux.HandleFunc("POST /like/thread", app.LikeHandler.LikeThread)
	mux.HandleFunc("POST /like/comment", app.LikeHandler.LikeComment)

	// Image routes
	mux.HandleFunc("POST /image", app.ImageHandler.UploadImage)

	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("GET /web/static/", http.StripPrefix("/web/static/", fs))

	return mux

}
