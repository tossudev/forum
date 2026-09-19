package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/render"
	"forum/internal/routes"
	"forum/internal/seed"
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
	renderer := render.NewRenderer()
	handler := routes.InitHandlers(db, validator, renderer)

	err = setupDB(db, cfg.Reset, cfg.MigrationsPath)
	if err != nil {
		slog.Error("error setting up database", "err", err)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	slog.Info("starting server", "addr", server.Addr)
	slog.Error("server failure", "err", server.ListenAndServe())
}

func setupDB(db *sql.DB, reset bool, path string) error {
	if reset {
		//clear database and see with dummy data
		if err := seed.ResetDatabase(db, path); err != nil {
			slog.Error("resetting database", "err", err)
			return err
		}
	} else {
		//create database without any dummy data
		if err := database.Migrate(db, path); err != nil {
			slog.Error("migrating database", "err", err)
			return err
		}
	}
	return nil
}
