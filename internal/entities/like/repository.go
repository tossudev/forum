package like

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/errs"
)

type LikeRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) GetThreadLike(ctx context.Context, thread_id, user_id int) (*Like, error) {
	exists, err := r.CheckThreadExists(ctx, thread_id)
	if err != nil {
		return nil, fmt.Errorf("get thread like: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("get thread like: %w", errs.ErrNotFound)
	}

	query := "SELECT thread_like FROM thread_likes WHERE thread_id = ? AND user_id = ?;"
	var liked bool
	var like Like
	if err := r.db.QueryRowContext(ctx, query, thread_id, user_id).Scan(&liked); err != nil {
		if err == sql.ErrNoRows {
			return &like, nil // error is nil because non-existent like reaction is valid
		}
		return nil, fmt.Errorf("get thread like: %w", err)
	}

	like.Exists = true
	like.Liked = liked

	return &like, nil
}

func (r *LikeRepository) RemoveThreadLike(ctx context.Context, thread_id, user_id int) error {
	query := "DELETE FROM thread_likes WHERE thread_id = ? AND user_id = ?"
	result, err := r.db.ExecContext(ctx, query, thread_id, user_id)

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("remove thread like: %w", err)
	}

	return nil
}

func (r *LikeRepository) UpdateThreadLike(ctx context.Context, req ThreadLikeRequest) error {
	query := "UPDATE thread_likes SET thread_like = ? WHERE thread_id = ? AND user_id = ?"
	if _, err := r.db.ExecContext(ctx, query, req.Like, req.ThreadID, req.UserID); err != nil {
		return fmt.Errorf("update thread like: %w", err)
	}

	return nil
}

// CheckThreadExists returns true if the thread with the given id exists
func (r *LikeRepository) CheckThreadExists(ctx context.Context, threadID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM threads WHERE id = ?);`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, threadID).Scan(&exists)
	return exists, err
}

// CheckCommentExists returns true if the thread with the given id exists
func (r *LikeRepository) CheckCommentExists(ctx context.Context, commentID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM comments WHERE id = ?);`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, commentID).Scan(&exists)
	return exists, err
}

func (r *LikeRepository) LikeThread(ctx context.Context, req ThreadLikeRequest) error {
	query := "INSERT INTO thread_likes (thread_id, user_id, thread_like) VALUES (?, ?, ?)"
	if _, err := r.db.ExecContext(ctx, query, req.ThreadID, req.UserID, req.Like); err != nil {
		return fmt.Errorf("like thread: %w", err)
	}

	return nil
}

func (r *LikeRepository) GetCommentLike(ctx context.Context, comment_id, user_id int) (*Like, error) {
	exists, err := r.CheckCommentExists(ctx, comment_id)
	if err != nil {
		return nil, fmt.Errorf("get comment like: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("get comment like: %w", errs.ErrNotFound)
	}

	query := "SELECT comment_like FROM comment_likes WHERE comment_id = ? AND user_id = ?;"
	var liked bool
	var like Like

	if err := r.db.QueryRowContext(ctx, query, comment_id, user_id).Scan(&liked); err != nil {
		if err == sql.ErrNoRows {
			return &like, nil // error is nil because non-existent like reaction is valid
		}
		return nil, fmt.Errorf("get comment like: %w", err)
	}

	like.Exists = true
	like.Liked = liked

	return &like, nil
}

func (r *LikeRepository) RemoveCommentLike(ctx context.Context, comment_id, user_id int) error {
	query := "DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?"
	result, err := r.db.ExecContext(ctx, query, comment_id, user_id)

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("remove comment like: %w", err)
	}

	return nil
}

func (r *LikeRepository) UpdateCommentLike(ctx context.Context, req CommentLikeRequest) error {
	query := "UPDATE comment_likes SET comment_like = ? WHERE comment_id = ? AND user_id = ?"
	if _, err := r.db.ExecContext(ctx, query, req.Like, req.CommentID, req.UserID); err != nil {
		return fmt.Errorf("update comment like: %w", err)
	}

	return nil
}

func (r *LikeRepository) LikeComment(ctx context.Context, req CommentLikeRequest) error {
	query := "INSERT INTO comment_likes (comment_id, user_id, comment_like) VALUES (?, ?, ?)"
	if _, err := r.db.ExecContext(ctx, query, req.CommentID, req.UserID, req.Like); err != nil {
		return fmt.Errorf("like comment: %w", err)
	}

	return nil
}
