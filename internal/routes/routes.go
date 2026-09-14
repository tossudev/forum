package routes

import (
	"net/http"
)

func GetRoutes(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /users", app.UserHandler.RegisterUser)

	// Thread routes
	mux.HandleFunc("GET /categories/{id}", app.ThreadHandler.GetByCategory)

	// Like routes
	mux.HandleFunc("POST /like/thread", app.LikeHandler.LikeThread)
	mux.HandleFunc("POST /like/comment", app.LikeHandler.LikeComment)

	return mux

}
