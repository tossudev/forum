package routes

import (
	"net/http"
)

func GetRoutes(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("Get /categories/{id}", app.ThreadHandler.GetByCategory)

	return mux
	
}
