package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ImageService struct {
	repo *ImageRepository
}

func NewService(r *ImageRepository) *ImageService {
	return &ImageService{repo: r}
}

func (s *ImageService) UploadImage(ctx context.Context, r *http.Request) error {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}
