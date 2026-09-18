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

	//"forum/internal/entities/comment"
	"forum/internal/entities/like"
	"forum/internal/password"
)

func ResetDatabase(db *sql.DB, path string) error {
	ctx := context.Background()

	query := `
	DROP TABLE IF EXISTS sessions;
	DROP TABLE IF EXISTS comments;
	DROP TABLE IF EXISTS thread_likes;
	DROP TABLE IF EXISTS comment_likes;
	DROP TABLE IF EXISTS threads;
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS categories;
	DROP TABLE IF EXISTS images;
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
	//CommentService *comment.CommentService
	LikeService *like.LikeService
}

func initSeedHandlers(db *sql.DB) SeedApp {
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)

	/*commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)*/

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)

	app := SeedApp{
		UserService:     userService,
		CategoryService: categoryService,
		ThreadService:   threadService,
		//CommentService: commentService,
		LikeService: likeService,
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
	/*if err := seedComments(ctx, app); err != nil {
		return fmt.Errorf("seedComments: %w", err)
	}*/
	if err := seedLikes(ctx, app); err != nil {
		return fmt.Errorf("seedLikes: %w", err)
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
			Email:    "imjustakid@email.com",
			Password: password.Password{},
			RoleID:   1,
		},
		{
			Username: "Marie",
			Email:    "mpp@thegoat.ca",
			Password: password.Password{},
			RoleID:   1,
		},
		{
			Username: "MontrealGoalie",
			Email:    "ard@montreal.ca",
			Password: password.Password{},
			RoleID:   1,
		},
		{
			Username: "",
			Email:    "",
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
		err := app.CategoryService.Create(ctx, &newCategory)
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
			CategoryID: 2,
		},
		{
			Title:      "The best book I've read this week",
			Body:       "It's called There Is No Antimimetic Division. #scifi",
			AuthorID:   4,
			CategoryID: 2,
		},
		{
			Title:      "My favourite author",
			Body:       "",
			AuthorID:   2,
			CategoryID: 4,
		},
		{
			Title:      "What are some good libraries in Helsinki???",
			Body:       "Help! I am visiting Helsinki for the weekend and would like to know, what are your favourite libraries in the city?? I will be there for three days and want to visit any cool or unique libraries you have. Also, any cool bookstores? Thanks!",
			AuthorID:   1,
			CategoryID: 1,
		},
		{
			Title: "The history of natural history",
			Body: `Let me tell you all about my favourite nonfiction genre. It started with Moby Dick (I know I know not nonfiction). When I was reading Moby Dick, 
			probably around the time where he was listing literally every kind of whale, I got curious if his whale stuff was based on any contemporary understanding of whales or if he was 
			just making stuff up. Eventually, while visiting London, I found a book called Leviathan by Philip Hoare. This book details Hoare's own obsession with whales through Moby Dick / Herman Mellville
			and through humanity's understanding of whales, and the sperm whale, through history. My enjoyment of this book led me to another book called Why Fish Don't Exist by Lulu Miller, which 
			is half personal memoir and half biography of a natural history dude (natural historialist?) from the 1800's who attempted to chronicle every fish in the world. Most recently, I have
			enjoyed a book called Beasts of the Sea (Iida Turpeinen), which details humanity's discovery, extermination of, and then retrospective fascination with the Stellar's sea cow. Beasts of the Sea
			is listed as fiction, but it feels like a history, which brings me back to Moby Dick...`,
			AuthorID:   1,
			CategoryID: 1,
		},
		{
			Title:      "",
			Body:       "",
			AuthorID:   1,
			CategoryID: 1,
		},
	}

	for _, newThread := range threads {
		err := app.ThreadService.Create(ctx, &newThread)
		if err != nil {
			return err
		}
	}

	return nil
}

/*func seedComments(ctx context.Context, app *SeedApp) error {
	comments := []comment.Comment{
		{
			Body: "What book?",
			ThreadID: 1,
			AuthorID: 1,
		},
		{
			Body: "What do you mean?",
			ThreadID: 1,
			AuthorID: 3,
		},
		{
			Body: "What was the book?",
			ThreadID: 1,
			AuthorID: 1,
		},
		{
			Body: "What book??",
			ThreadID: 1,
			AuthorID: 3,
		},
		{
			Body: "The name of the book that you read!!",
			ThreadID: 1,
			AuthorID: 1,
		},
		{
			Body: `I found Katabasis (in English) in a small airport bookstore in Croatia. Not sure how I feel about finding one of my favourite authors of the last few years
			in an airport...`,
			ThreadID: 3,
			AuthorID: 4,
		},
		{
			Body: "Maybe Croatians just have good taste *shrug*",
			ThreadID: 3,
			AuthorID: 2,
		},
		{
			Body: "You may be right. There was also a beautiful version of Wuthering Heights that I was tempted to pick up. No Dean Koontz in Croatia, I guess.",
			ThreadID: 3,
			AuthorID: 4,
		},
		{
			Body: "Do people read Dean Koontz in Europe?",
			ThreadID: 3,
			AuthorID: 5,
		},
		{
			Body: "",
			ThreadID: 1,
			AuthorID: 1,
		},
	}

	for _, newComment := range comments {
		_, err := app.CommentService.Create(ctx, &newComment)
		if err != nil {
			return err
		}
	}

	return nil
}*/

func seedLikes(ctx context.Context, app *SeedApp) error {
	threadLikes := []like.ThreadLikeRequest{
		{
			ThreadID: 1,
			UserID:   5,
			Like:     true,
		},
		{
			ThreadID: 1,
			UserID:   1,
			Like:     false,
		},
	}

	/*commentLikes := []like.CommentLikeRequest{
		{
			CommentID: 1,
			UserID: 1,
			Like:
		},
	}*/

	for _, newThreadLike := range threadLikes {
		err := app.LikeService.LikeThread(ctx, newThreadLike)
		if err != nil {
			return err
		}
	}

	/*for _, newCommentLike := range commentLikes {
		err := app.LikeService.LikeComment(ctx, newCommentLike)
		if err != nil {
			return err
		}
	}*/

	return nil
}
