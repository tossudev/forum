package like

import (
	"context"
	"database/sql"
	"fmt"
)

type LikeService struct {
	repo *LikeRepository
}

func NewService(r *LikeRepository) *LikeService {
	return &LikeService{repo: r}
}

func (s *LikeService) LikeThread(ctx context.Context, req ThreadLikeRequest) error {
	// Retrieve existing like/dislike
	like, err := s.repo.GetThreadLike(ctx, req.ThreadID, req.UserID)
	// If no like/dislike, apply like/dislike
	if err != nil {
		if err == sql.ErrNoRows {
			return s.repo.LikeThread(ctx, req)
		} else {
			return fmt.Errorf("get thread like: %w", err)
		}
	}

	// Swap between like/dislike
	if like != req.Like {
		return s.repo.UpdateThreadLike(ctx, req)
	}

	// Double positive, cancel out
	return s.repo.RemoveThreadLike(ctx, req.ThreadID, req.UserID)
}

func (s *LikeService) LikeComment(ctx context.Context, req CommentLikeRequest) error {
	// Retrieve existing like/dislike
	like, err := s.repo.GetCommentLike(ctx, req.CommentID, req.UserID)
	// If no like/dislike, apply like/dislike
	if err != nil {
		if err == sql.ErrNoRows {
			return s.repo.LikeComment(ctx, req)
		} else {
			return fmt.Errorf("get comment like: %w", err)
		}
	}

	// Swap between like/dislike
	if like != req.Like {
		return s.repo.UpdateCommentLike(ctx, req)
	}

	// Double positive, cancel out
	return s.repo.RemoveCommentLike(ctx, req.CommentID, req.UserID)
}
