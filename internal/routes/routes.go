package routes

import (
	"net/http"
)

func GetRoutes(app *App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /categories/{id}", app.ThreadHandler.GetByCategory)
	mux.HandleFunc("POST /threads", app.ThreadHandler.Create)
	mux.HandleFunc("GET /threads/{id}", app.ThreadHandler.GetByID)

	return mux

}
