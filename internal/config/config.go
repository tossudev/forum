package config

import "flag"

type Config struct {
	Port           int
	DbPath         string
	MigrationsPath string
}

func Load() Config {
	cfg := Config{MigrationsPath: "migrations/001_init.sql"}

	flag.IntVar(&cfg.Port, "port", 8080, "Server port number")
	flag.StringVar(&cfg.DbPath, "db", "forum.db", "Database file path")
	flag.Parse()

	return cfg
}
