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

	"forum/internal/entities/comment"
	"forum/internal/entities/like"
	"forum/internal/password"
)

func ResetDatabase(db *sql.DB, path string) error {
	ctx := context.Background()

	query := `
	DROP TABLE IF EXISTS sessions;
	DROP TABLE IF EXISTS thread_likes;
	DROP TABLE IF EXISTS comment_likes;
	DROP TABLE IF EXISTS images;
	DROP TABLE IF EXISTS comments;
	DROP TABLE IF EXISTS threads;
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS categories;
	DROP TABLE IF EXISTS threads_fts;
	DROP TABLE IF EXISTS comments_fts;
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
	CommentService  *comment.CommentService
	LikeService     *like.LikeService
}

func initSeedHandlers(db *sql.DB) SeedApp {
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)

	threadRepo := thread.NewRepository(db)
	threadService := thread.NewService(threadRepo)

	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)

	likeRepo := like.NewRepository(db)
	likeService := like.NewService(likeRepo)

	app := SeedApp{
		UserService:     userService,
		CategoryService: categoryService,
		ThreadService:   threadService,
		CommentService:  commentService,
		LikeService:     likeService,
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
	if err := seedComments(ctx, app); err != nil {
		return fmt.Errorf("seedComments: %w", err)
	}
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
			Avatar:   "avatar_e.png",
		},
		{
			Username: "JustinV",
			Email:    "v@justin.com",
			Password: password.Password{},
			RoleID:   1,
			Avatar:   "avatar_v.png",
		},
		{
			Username: "KevinMac2004",
			Email:    "imjustakid@email.com",
			Password: password.Password{},
			RoleID:   1,
			Avatar:   "avatar_s.png",
		},
		{
			Username: "Marie",
			Email:    "mpp@thegoat.ca",
			Password: password.Password{},
			RoleID:   1,
			Avatar:   "bee.gif",
		},
		{
			Username: "MontrealGoalie",
			Email:    "ard@montreal.ca",
			Password: password.Password{},
			RoleID:   1,
			Avatar:   "",
		},
	}

	for _, newUser := range users {
		password := "password"
		_, err := app.UserService.RegisterUser(ctx, &newUser, password)
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
			Body:       "I super loved this book!\nI read it in two days!!",
			AuthorID:   1,
			CategoryID: 2,
		},
		{
			Title:      "The best book I've read this week",
			Body:       "It's called There Is No Antimimetic Division.\n\n#scifi",
			AuthorID:   4,
			CategoryID: 2,
		},
		{
			Title:      "My favourite author",
			Body:       "Is... ?",
			AuthorID:   2,
			CategoryID: 4,
		},
		{
			Title:      "What are some good libraries in Helsinki???",
			Body:       "Help! I am visiting Helsinki for the weekend and would like to know, what are your favourite libraries in the city??\nI will be there for three days and want to visit any cool or unique libraries you have. Also, any cool bookstores? Thaaanks!",
			AuthorID:   1,
			CategoryID: 1,
		},
		{
			Title:      "The history of natural history",
			Body:       `Let me tell you all about my favourite nonfiction genre. It started with Moby Dick (I know I know not nonfiction). When I was reading Moby Dick, probably around the time where he was listing literally every kind of whale, I got curious if his whale stuff was based on any contemporary understanding of whales or if he was just making stuff up. Eventually, while visiting London, I found a book called Leviathan by Philip Hoare. This book details Hoare's own obsession with whales through Moby Dick / Herman Mellville and through humanity's understanding of whales, and the sperm whale, through history. My enjoyment of this book led me to another book called Why Fish Don't Exist by Lulu Miller, which is half personal memoir and half biography of a natural history dude (natural historialist?) from the 1800's who attempted to chronicle every fish in the world. Most recently, I have enjoyed a book called Beasts of the Sea (Iida Turpeinen), which details humanity's discovery, extermination of, and then retrospective fascination with the Stellar's sea cow. Beasts of the Sea is listed as fiction, but it feels like a history, which brings me back to Moby Dick...`,
			AuthorID:   3,
			CategoryID: 1,
		},
		{
			Title:      "Cool GIFs!",
			Body:       "Comment your favourite GIFs!",
			AuthorID:   4,
			CategoryID: 5,
		},
		{
			Title:      "Reading in Finnish",
			Body:       "I have been wanting to try to read some scifi/fantasy books originally written in Finnish. I have read some Finnish books that were translated into Finnish. I read Turtles All The Way Down in Finnish. And I have read a couple of the Muumin books (though I always forget those were originall written in Swedish, not Finnish).\nHave any of you done any reading in Finnish? What do you recommend?",
			AuthorID:   4,
			CategoryID: 2,
		},
		{
			Title:      "Haruki Murakami",
			Body:       "Hardboiled Wonderland and the End of the World was good, but in the other books I've tried to read I find myself frustrated with his passive protagonists that just sit around and let their lives happen to them.",
			AuthorID:   2,
			CategoryID: 4,
		},
		{
			Title:      "Onko tällä ketään joka puhuu suomea",
			Body:       "Olisi kiva puhua",
			AuthorID:   4,
			CategoryID: 5,
		},
		{
			Title:      "📚📚📚",
			Body:       "📖📕📘📙🤓\n🦁🦁🦁\nBOOKS",
			AuthorID:   1,
			CategoryID: 2,
		},
		{
			Title:      "Nonfiction books",
			Body:       "Looking for some interesting nonfiction to read (nothing about whales, Kevin!)",
			AuthorID:   2,
			CategoryID: 2,
		},
		{
			Title:      "Books that have been made into movies",
			Body:       "What are your favourite books that have been made into movies? Leave your recommendations in the comments!",
			AuthorID:   1,
			CategoryID: 2,
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

func seedComments(ctx context.Context, app *SeedApp) error {
	comments := []comment.Comment{
		{
			Body:     "What book?",
			ThreadID: 1,
			AuthorID: 3,
		},
		{
			Body:     "What do you mean?",
			ThreadID: 1,
			AuthorID: 1,
		},
		{
			Body:     "What was the book?",
			ThreadID: 1,
			AuthorID: 3,
		},
		{
			Body:     "What book??",
			ThreadID: 1,
			AuthorID: 1,
		},
		{
			Body:     "The name of the book that you read!!",
			ThreadID: 1,
			AuthorID: 3,
		},
		{
			Body:     "Is it a book?",
			ThreadID: 1,
			AuthorID: 1,
		},
		{
			Body:     "What kind of book?",
			ThreadID: 1,
			AuthorID: 5,
		},
		{
			Body:     "What is your favourite book?",
			ThreadID: 1,
			AuthorID: 4,
		},
		{
			Body:     "What kind of book is it?",
			ThreadID: 2,
			AuthorID: 4,
		},
		{
			Body:     "Tell me about the book",
			ThreadID: 1,
			AuthorID: 4,
		},
		{
			Body:     "What kind of books do they write?",
			ThreadID: 3,
			AuthorID: 1,
		},
		{
			Body:     "The best kinds of books",
			ThreadID: 3,
			AuthorID: 2,
		},
		{
			Body: `I found Katabasis (in English) in a small airport bookstore in Croatia. Not sure how I feel about finding one of my favourite authors of the last few years
			in an airport...`,
			ThreadID: 3,
			AuthorID: 4,
		},
		{
			Body:     "Maybe Croatians just have good taste *shrug*",
			ThreadID: 3,
			AuthorID: 2,
		},
		{
			Body:     "You may be right. There was also a beautiful version of Wuthering Heights that I was tempted to pick up. No Dean Koontz in Croatia, I guess.",
			ThreadID: 3,
			AuthorID: 4,
		},
		{
			Body:     "Do people read Dean Koontz in Europe?",
			ThreadID: 3,
			AuthorID: 5,
		},
		{
			Body:     `<img src="https://hadyn.org/mediagallery/gifs/m_catcamp.gif">`,
			ThreadID: 6,
			AuthorID: 5,
		},
		{
			Body:     `<img src="https://coolgifs.neocities.org/gifs/137.gif">`,
			ThreadID: 6,
			AuthorID: 5,
		},
		{
			Body: `I LOVE CATS!
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20513.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20292.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20293.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20725.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20726.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20727.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20728.gif">
			<img src="https://bettysgraphics.neocities.org/images/animals/cat%20729.gif">
			`,
			ThreadID: 6,
			AuthorID: 5,
		},
		{
			Body:     "Aww, man. No fair!",
			ThreadID: 11,
			AuthorID: 3,
		},
		{
			Body:     `<img src="https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExZ253em9rYTZsODh2YTljdzc5cGlncXdzNzNoaWRsanl6aG9jMTM3YiZlcD12MV9naWZzX3NlYXJjaCZjdD1n/FaKV1cVKlVRxC/giphy.gif">`,
			ThreadID: 11,
			AuthorID: 4,
		},
		{
			Body:     `<img src="https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExcjF1Zm92cHNoMWNnc3ZxZ3JwcGNncGFrZzQ2dXd6M3lvcjdhNnVzaSZlcD12MV9naWZzX3NlYXJjaCZjdD1n/4EiGNSTfy4WC4/giphy.gif">`,
			ThreadID: 6,
			AuthorID: 4,
		},
		{
			Body:     "Suomi mainittu!",
			ThreadID: 9,
			AuthorID: 5,
		},
		{
			Body:     "Torille!🇫🇮",
			ThreadID: 9,
			AuthorID: 1,
		},
		{
			Body:     "Torille!🇫🇮🇫🇮🇫🇮",
			ThreadID: 9,
			AuthorID: 3,
		},
		{
			Body:     "What?",
			ThreadID: 9,
			AuthorID: 4,
		},
		{
			Body:     "Guys... this is a Finnish forum...",
			ThreadID: 9,
			AuthorID: 2,
		},
		{
			Body:     `<img src="https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExbzhsODYyM3pmYmpkNmRlYWgwMWY2N3hqOWVuN3hnb3huY3c4dm1ndyZlcD12MV9naWZzX3JlbGF0ZWQmY3Q9Zw/p4gI9UwAoTs75VMwPA/giphy.gif">`,
			ThreadID: 9,
			AuthorID: 3,
		},
	}

	for _, newComment := range comments {
		_, _, err := app.CommentService.Create(ctx, &newComment)
		if err != nil {
			return err
		}
	}

	return nil
}

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

	commentLikes := []like.CommentLikeRequest{
		{
			CommentID: 1,
			UserID:    1,
			Like:      true,
		},
		{
			CommentID: 1,
			UserID:    2,
			Like:      false,
		},
		{
			CommentID: 1,
			UserID:    3,
			Like:      true,
		},
		{
			CommentID: 21,
			UserID:    1,
			Like:      true,
		},
		{
			CommentID: 21,
			UserID:    2,
			Like:      false,
		},
		{
			CommentID: 21,
			UserID:    3,
			Like:      true,
		},
		{
			CommentID: 21,
			UserID:    4,
			Like:      true,
		},
		{
			CommentID: 21,
			UserID:    5,
			Like:      true,
		},
		{
			CommentID: 20,
			UserID:    3,
			Like:      true,
		},
		{
			CommentID: 26,
			UserID:    5,
			Like:      true,
		},
		{
			CommentID: 26,
			UserID:    1,
			Like:      true,
		},
		{
			CommentID: 26,
			UserID:    3,
			Like:      true,
		},
		{
			CommentID: 28,
			UserID:    5,
			Like:      true,
		},
		{
			CommentID: 28,
			UserID:    1,
			Like:      true,
		},
		{
			CommentID: 28,
			UserID:    3,
			Like:      true,
		},
	}

	for _, newThreadLike := range threadLikes {
		err := app.LikeService.LikeThread(ctx, newThreadLike)
		if err != nil {
			return err
		}
	}

	for _, newCommentLike := range commentLikes {
		err := app.LikeService.LikeComment(ctx, newCommentLike)
		if err != nil {
			return err
		}
	}

	return nil
}
