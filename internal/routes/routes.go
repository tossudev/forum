package routes

import (
	"forum/internal/session"
	"net/http"
)

func GetRoutes(app *App, sm *session.SessionManager) *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("GET /users/register", app.UserHandler.RegisterPage)
	mux.HandleFunc("GET /users/login", app.UserHandler.LoginPage)
	mux.HandleFunc("GET /profile", app.UserHandler.GetProfile)
	mux.HandleFunc("GET /profile/threads", app.UserHandler.ProfileThreadsPage)
	mux.HandleFunc("POST /users/register", app.UserHandler.RegisterUser)
	mux.HandleFunc("POST /users/login", app.UserHandler.Login)
	mux.HandleFunc("POST /users/logout", app.UserHandler.Logout)

	// Category routes
	mux.HandleFunc("POST /categories", app.CategoryHandler.Create)
	mux.HandleFunc("GET /", app.CategoryHandler.GetAll)

	// Thread routes
	mux.HandleFunc("GET /categories/{id}", app.ThreadHandler.GetByCategory)
	mux.HandleFunc("GET /categories/{id}/create", app.ThreadHandler.CreateThreadPage)
	mux.HandleFunc("POST /categories/{id}/create", app.ThreadHandler.Create)
	mux.HandleFunc("GET /threads/{id}", app.ThreadHandler.DisplayThreadPage)
	mux.HandleFunc("DELETE /threads/{id}", app.ThreadHandler.Delete)

	//Comment routes
	mux.HandleFunc("GET /comments/{id}", app.CommentHandler.GetByID)
	mux.HandleFunc("GET /threads/{id}/comments", app.CommentHandler.GetByThread)
	mux.HandleFunc("POST /threads/{id}/comments", app.CommentHandler.Create)
	mux.HandleFunc("DELETE /comments/{id}", app.CommentHandler.Delete)

	// Like routes
	mux.HandleFunc("POST /like/thread", app.LikeHandler.LikeThread)
	mux.HandleFunc("POST /like/comment", app.LikeHandler.LikeComment)

	// Image routes
	mux.HandleFunc("POST /image", app.ImageHandler.UploadImage)

	//Search routes
	mux.HandleFunc("GET /search", app.SearchHandler.Search)

	// TODO: consider using embedding for static
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("GET /web/static/", http.StripPrefix("/web/static/", fs))

	fsUploads := http.FileServer(http.Dir("./uploads/"))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fsUploads))

	return mux

}
