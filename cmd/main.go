package main

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/http" //package name???
	"log/slog"
	"net/http"
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

	validate := http.initValidator() //package name???

	app := http.initApp(db, validate) //package name???

	if err := database.Migrate(db, cfg.MigrationsPath); err != nil {
		slog.Error("migrating database", "err", err)
		return
	}

	mux := http.NewServeMux()
	http.GetRoutes(mux, app) //package name???

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	slog.Info("starting server", "addr", server.Addr)
	slog.Error("server failure", "err", server.ListenAndServe())
}
