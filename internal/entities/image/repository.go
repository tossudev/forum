package image

import (
	"database/sql"
)

type ImageRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) UploadImage() {

}
