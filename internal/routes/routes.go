package routes

import (
	"forum/internal/middleware"

)

func GetRoutes(app *App) *ServeMux{
	mux := http.NewServeMux()

	mux.HandlFunc("Get /categories/{id}", app.ThreadHandler.GetThreadsByCategory)


	// Middleware
	handler := middleware.Logger(mux)
	handler = middleware.Recovery(handler)

	return handler
	
}