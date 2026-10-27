package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

func AddCategories(db *sql.DB) error {
	ctx := context.Background()

	var exists bool
	queryExists := `SELECT EXISTS (SELECT 1 FROM categories)`

	if err := db.QueryRowContext(ctx, queryExists).Scan(&exists); err != nil {
		return fmt.Errorf("prep db, check if categories exist: %w", err)
	}

	if !exists {
		if err := InitCategories(ctx, db); err != nil {
			return fmt.Errorf("prep db, initialise empty categories: %w", err)
		}

	}
	return nil
}

func InitCategories(ctx context.Context, db *sql.DB) error {
	app := initSeedHandlers(db)

	if err := seedCategories(ctx, &app); err != nil {
		return err
	}
	slog.Info("Categories added to fresh database")
	return nil
}
