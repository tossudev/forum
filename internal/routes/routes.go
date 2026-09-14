package routes

import (
	"net/http"
)

func GetRoutes(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("POST /users", app.UserHandler.RegisterUser)

	// Thread routes
	mux.HandleFunc("Get /categories/{id}", app.ThreadHandler.GetByCategory)

	return mux

}
