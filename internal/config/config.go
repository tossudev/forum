package config

import (
	"flag"
	"log/slog"
)

type Config struct {
	Port           int
	DbPath         string
	MigrationsPath string
	Reset          bool
	LogLevel       slog.Level
}

const TimeFormat string = "20060102T150405"

func Load() Config {
	cfg := Config{
		MigrationsPath: "migrations/001_init.sql",
	}

	flag.IntVar(&cfg.Port, "port", 8080, "Server port number")
	flag.StringVar(&cfg.DbPath, "db", "forum.db", "Database file path")
	flag.BoolVar(&cfg.Reset, "reset", false, "Reset the database")
	flag.TextVar(&cfg.LogLevel, "log-level", slog.LevelDebug, "Log level: debug, info, warn, or error")
	flag.Parse()

	return cfg
}
