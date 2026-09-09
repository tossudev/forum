package main

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
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

	if err := database.Migrate(db, cfg.MigrationsPath); err != nil {
		slog.Error("migrating database", "err", err)
		return
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	slog.Info("starting server", "addr", server.Addr)
	slog.Error("server failure", "err", server.ListenAndServe())
}
