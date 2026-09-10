package http

import (
	"net/http"
	"forum/internal/thread"
)

func GetRoutes(mux *http.ServeMux, app *App) {

	mux.HandlFunc("Get /categories/{id}", app.GetThreadsByCategory)
	
}