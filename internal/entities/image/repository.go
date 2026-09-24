package image

import (
	"context"
	"database/sql"
	"fmt"
)

type ImageRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Create(ctx context.Context, path string, req ImageRequest) error {
	query := "INSERT INTO images (image_path, thread_id, comment_id) VALUES (?, ?, ?)"
	if _, err := r.db.ExecContext(ctx, query, path, req.ThreadID, req.CommentID); err != nil {
		return fmt.Errorf("create image err: %w", err)
	}

	return nil
}
