package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"forum/internal/database"
	"forum/internal/entities/category"
	"forum/internal/entities/thread"
	"forum/internal/entities/user"
	"forum/internal/password"
)

func ResetDatabase(db *sql.DB, path string) error {
	ctx := context.Background()

	query := `
	DROP TABLE IF EXISTS threads;
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS comments;
	DROP TABLE IF EXISTS thread_likes;
	DROP TABLE IF EXISTS comment_likes;
	DROP TABLE IF EXISTS categories;
	DROP TABLE IF EXISTS images;
	DROP TABLE IF EXISTS sessions;
	`

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	slog.Info("Deleted all data from database")

	if err := database.Migrate(db, path); err != nil {
		return err
	}
	slog.Info("Reset database schema")

	app := initSeedHandlers(db)

	if err := seedDatabase(ctx, &app); err != nil {
		return err
	}

	return nil
}

type SeedApp struct {
	UserService     *user.UserService
	CategoryService *category.CategoryService
	ThreadService   *thread.ThreadService
}

func initSeedHandlers(db *sql.DB) SeedApp {
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)

	app := SeedApp{
		UserService:     userService,
		CategoryService: categoryService,
		ThreadService:   threadService,
	}

	return app
}

func seedDatabase(ctx context.Context, app *SeedApp) error {
	if err := seedUsers(ctx, app); err != nil {
		return fmt.Errorf("seedUsers: %w", err)
	}
	if err := seedCategories(ctx, app); err != nil {
		return fmt.Errorf("seedCategories: %w", err)
	}
	if err := seedThreads(ctx, app); err != nil {
		return fmt.Errorf("seedThreads: %w", err)
	}
	return nil
}

func seedUsers(ctx context.Context, app *SeedApp) error {

	users := []user.User{
		{
			Username: "Name",
			Email:    "name@name.com",
			Password: password.Password{},
			RoleID:   1,
		},
		{
			Username: "JustinV",
			Email:    "v@justin.com",
			Password: password.Password{},
			RoleID:   1,
		},
		{
			Username: "KevinMac2004",
			Email:    "KevinMac@email.com",
			Password: password.Password{},
			RoleID:   1,
		},
		{
			Username: "Marie",
			Email:    "mpp@thegoat.ca",
			Password: password.Password{},
			RoleID:   1,
		},
	}

	for _, newUser := range users {
		password := "password"
		err := app.UserService.RegisterUser(ctx, &newUser, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedCategories(ctx context.Context, app *SeedApp) error {
	categories := []category.Category{
		{
			Name: "General",
		},
		{
			Name: "Books",
		},
		{
			Name: "Genres",
		},
		{
			Name: "Authors",
		},
		{
			Name: "Off Topic",
		},
	}

	for _, newCategory := range categories {
		_, err := app.CategoryService.Create(ctx, &newCategory)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedThreads(ctx context.Context, app *SeedApp) error {

	threads := []thread.Thread{
		{
			Title:      "I loved this book",
			Body:       "I super loved this book! I read it in two days!!",
			AuthorID:   3,
			CategoryID: 1,
		},
	}

	for _, newThread := range threads {
		_, err := app.ThreadService.Create(ctx, &newThread)
		if err != nil {
			return err
		}
	}

	return nil
}
