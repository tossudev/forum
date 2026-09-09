package main

import (
	"forum/internal/config"
	"log/slog"
)

func main() {
	cfg := config.Load()
	slog.Info("config loaded:", "config", cfg)
}
