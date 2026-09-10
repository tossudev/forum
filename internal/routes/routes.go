package routes

import (

)

func GetRoutes(mux *http.ServeMux, app *App) {

	mux.HandlFunc("Get /categories/{id}", app.ThreadHandler.GetThreadsByCategory)
	
}