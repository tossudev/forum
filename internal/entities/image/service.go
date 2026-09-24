package image

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"uuid"

	"forum/internal/config"
	"forum/internal/errs"
)

var validTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/jpg":  {},
	"image/png":  {},
	"image/gif":  {},
}

type ImageService struct {
	repo *ImageRepository
}

func NewService(r *ImageRepository) *ImageService {
	return &ImageService{repo: r}
}

func (s *ImageService) UploadImage(ctx context.Context, r *http.Request, req ImageRequest) (string, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	valid, err := validateImage(file)
	if err != nil {
		return "", fmt.Errorf("image validation failed: %w", err)
	}

	if !valid {
		return "", errs.ErrInvalidFiletype
	}

	imageName := uuid.NewV7().String() + filepath.Ext(fileHeader.Filename)
	imagePath := filepath.Join(config.UploadsPath, imageName)
	if err := createFile(file, imagePath); err != nil {
		return "", fmt.Errorf("image creation failed: %w", err)
	}

	// If database insertion fails, remove the image from disk
	if err := s.repo.Create(ctx, imageName, req); err != nil {
		_ = os.Remove(imagePath)
		return "", fmt.Errorf("database image creation failed: %w", err)
	}

	return imagePath, nil
}

func validateImage(file multipart.File) (bool, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	filetype := http.DetectContentType(buffer)

	if _, ok := validTypes[filetype]; !ok {
		return false, nil
	}

	return true, nil
}

func createFile(file multipart.File, imagePath string) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	dst, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err := dst.Write(fileBytes); err != nil {
		return err
	}

	return nil
}
