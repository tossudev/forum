package thread

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"forum/internal/errs"
	"forum/internal/pagination"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) GetByCategory(ctx context.Context, id int, pagination pagination.Pagination) ([]Thread, error) {
	var threads []Thread

	query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE category_id = ? ORDER BY date_created ASC LIMIT ? OFFSET ?;"

	rows, err := r.db.QueryContext(ctx, query, id, pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID); err != nil {
			return nil, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		queryCount := "SELECT COUNT(*) FROM comments WHERE thread_id = ?"
		if err := r.db.QueryRowContext(ctx, queryCount, thread.ID).Scan(&thread.ReplyCount); err != nil {
			return nil, fmt.Errorf("Thread GetByCategory: %w", err)
		}

		threads = append(threads, thread)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Thread GetByCategory: %w", err)
	}

	return threads, nil
}

func (r *ThreadRepository) GetByID(ctx context.Context, id int) (*Thread, error) {
	query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE id = ?;"
	var thread Thread
	err := r.db.QueryRowContext(ctx, query, id).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("Thread GetByID: %w", err)
	}

	return &thread, nil

}

func (r *ThreadRepository) Create(ctx context.Context, thread *Thread) (int, error) {
	query := "INSERT INTO threads (title, body, date_created, author_id, category_id) VALUES (?, ?, ?, ?, ?);"
	result, err := r.db.ExecContext(ctx, query, thread.Title, thread.Body, thread.DateCreated, thread.AuthorID, thread.CategoryID)
	if err != nil {
		return 0, fmt.Errorf("Thread Create: %w", err)
	}

	threadID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Thread Create: getting last insert id: %w", err)
	}

	return int(threadID), nil
}

func (r *ThreadRepository) Delete(ctx context.Context, threadID int, userID int) error {
	thread, err := r.GetByID(ctx, threadID)
	if err != nil {
		return fmt.Errorf("Thread search: %w", err)
	}

	if userID != thread.AuthorID {
		return fmt.Errorf("%w: user is not thread creator", errs.ErrUnauthorized)
	}

	query := "DELETE from threads WHERE id = ?;"
	_, err = r.db.ExecContext(ctx, query, threadID)
	if err != nil {
		return fmt.Errorf("Thread delete: %w", err)
	}

	return nil
}

func (r *ThreadRepository) DisplayThreadPage(ctx context.Context, threadID int, username string) (*ThreadData, error) {
	threadQuery := `SELECT threads.id, title, body, threads.date_created, threads.author_id, threads.category_id, username FROM users INNER JOIN threads ON users.id = threads.author_id WHERE threads.id = ?;`
	//query := `SELECT threads.id, title, body, date_created, username, thread_like, thread_likes.user_id FROM threads INNER JOIN users ON users.id = threads.author_id LEFT JOIN thread_likes ON threads.id = thread_likes.thread_id;`
	//query := "SELECT id, title, body, date_created, author_id, category_id FROM threads WHERE id = ?;"

	var authorName string
	var thread Thread

	err := r.db.QueryRowContext(ctx, threadQuery, threadID).Scan(&thread.ID, &thread.Title, &thread.Body, &thread.DateCreated, &thread.AuthorID, &thread.CategoryID, &authorName)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	//this is a separate query because otherwise there would be a distinct row in the above query for every like on every thread
	threadLikesQuery := `SELECT username, thread_like FROM thread_likes INNER JOIN users ON users.id = thread_likes.user_id WHERE thread_id = ?`

	rows, err := r.db.QueryContext(ctx, threadLikesQuery, threadID)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}
	var likeUsers []string
	var dislikeUsers []string
	for rows.Next() {
		var likeUser string
		var like bool
		if err := rows.Scan(&likeUser, &like); err != nil {
			return nil, fmt.Errorf("DisplayThreadPage: %w", err)
		}
		if like {
			likeUsers = append(likeUsers, likeUser)
		} else {
			dislikeUsers = append(dislikeUsers, likeUser)
		}
	}

	//get all comments for this thread plus username of comment author
	threadCommentsQuery := `SELECT comments.id, comments.body, comments.thread_id, username, comments.date_created FROM comments INNER JOIN users ON comments.author_id = users.id WHERE comments.thread_id = ?;`
	commentRows, err := r.db.QueryContext(ctx, threadCommentsQuery, threadID)
	if err != nil {
		return nil, fmt.Errorf("DisplayThreadPage: %w", err)
	}

	//query for comment likes/dislikes
	commentLikesQuery := `SELECT username, comment_like FROM comment_likes INNER JOIN users ON users.id = comment_likes.user_id WHERE comment_id = ?`

	var comments []CommentData

	//for every comment on thread:
	for commentRows.Next() {
		var comment CommentData
		if err := commentRows.Scan(&comment.ID, &comment.Body, &comment.ThreadID, &comment.AuthorName, &comment.DateCreated); err != nil {
			return nil, fmt.Errorf("DisplayThreadPage: %w", err)
		}

		//get comment likes/dislikes
		likeRows, err := r.db.QueryContext(ctx, commentLikesQuery, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("DispleyThreadPage: %w", err)
		}
		var commentLikeUsers []string
		var commentDislikeUsers []string
		for likeRows.Next() {
			var commentLikeUser string
			var commentLike bool
			if err := likeRows.Scan(&commentLikeUser, &commentLike); err != nil {
				return nil, fmt.Errorf("DisplayThreadPage: %w", err)
			}
			if commentLike {
				commentLikeUsers = append(commentLikeUsers, commentLikeUser)
			} else {
				commentDislikeUsers = append(commentDislikeUsers, commentLikeUser)
			}
		}

		comment.NumLikes = len(commentLikeUsers)
		comment.LikeUsers = commentLikeUsers
		comment.NumDislikes = len(commentDislikeUsers)
		comment.DislikeUsers = commentDislikeUsers

		if len(username) > 0 && slices.Contains(comment.LikeUsers, username) {
			comment.LoggedUserLiked = true
		}
		if len(username) > 0 && slices.Contains(comment.DislikeUsers, username) {
			comment.LoggedUserDisliked = true
		}

		//append comment to comments
		comments = append(comments, comment)
	}

	data := ThreadData{
		Thread:       thread,
		AuthorName:   authorName,
		NumLikes:     len(likeUsers),
		LikeUsers:    likeUsers,
		NumDislikes:  len(dislikeUsers),
		DislikeUsers: dislikeUsers,
		Comments:     comments,
	}

	return &data, nil
}
