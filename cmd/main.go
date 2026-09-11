package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/routes" 
	"forum/internal/validate" 
)

func main() {
	cfg := config.Load()
	slog.Info("config loaded:", "config", cfg)

	db, err := database.Open(cfg.DbPath)
	if err != nil {
		slog.Error("opening database", "err", err)
		return
	}
	defer db.Close()
	slog.Info("connected to database", "db", cfg.DbPath)

	validator := validate.InitValidator() 
	handler := routes.InitHandlers(db, validator)  //and here

	if err := database.Migrate(db, cfg.MigrationsPath); err != nil {
		slog.Error("migrating database", "err", err)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	slog.Info("starting server", "addr", server.Addr)
	slog.Error("server failure", "err", server.ListenAndServe())
}
