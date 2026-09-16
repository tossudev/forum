package seed

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/database"
	"forum/internal/entities/user"
	"forum/internal/password"
)

func ResetDatabase(db *sql.DB, path string) error {
	ctx := context.Background()

	query := `
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS threads;
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
	fmt.Println("Deleted all data from database")

	if err := database.Migrate(db, path); err != nil {
		return err
	}
	fmt.Println("Reset database schema")

	app := initSeedHandlers(db)

	if err := seedDatabase(ctx, &app); err != nil {
		return err
	}

	return nil
}

type SeedApp struct {
	UserService *user.UserService
	//ThreadService *thread.ThreadService
}

func initSeedHandlers(db *sql.DB) SeedApp {
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)

	app := SeedApp{
		UserService: userService,
		//ThreadHandler: threadHandler,
	}

	return app
}

func seedDatabase(ctx context.Context, app *SeedApp) error {
	if err := seedUsers(ctx, app); err != nil {
		return fmt.Errorf("seedUsers: %w", err)
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
