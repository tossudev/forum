package like

import (
	"context"
	"database/sql"
	"fmt"
)

type LikeRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) GetThreadLike(ctx context.Context, thread_id, user_id int) (bool, bool, error) {
	var like bool
	query := "SELECT thread_like FROM thread_likes WHERE thread_id = ? AND user_id = ?;"

	row := r.db.QueryRow(query, thread_id, user_id)

	if err := row.Scan(&like); err != nil {
		if err == sql.ErrNoRows {
			return like, false, nil
		}
		return like, false, fmt.Errorf("get thread like: %w", err)
	}

	return like, true, nil
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

func (r *LikeRepository) LikeThread(ctx context.Context, req ThreadLikeRequest) error {
	query := "INSERT INTO thread_likes (thread_id, user_id, thread_like) VALUES (?, ?, ?)"
	if _, err := r.db.ExecContext(ctx, query, req.ThreadID, req.UserID, req.Like); err != nil {
		return fmt.Errorf("like thread: %w", err)
	}

	return nil
}

func (r *LikeRepository) GetCommentLike(ctx context.Context, comment_id, user_id int) (bool, bool, error) {
	var like bool
	query := "SELECT comment_like FROM comment_likes WHERE comment_id = ? AND user_id = ?;"

	row := r.db.QueryRow(query, comment_id, user_id)

	if err := row.Scan(&like); err != nil {
		if err == sql.ErrNoRows {
			return like, false, nil
		}
		return like, false, fmt.Errorf("get comment like: %w", err)
	}

	return like, true, nil
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
