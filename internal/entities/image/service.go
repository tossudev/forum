package image

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"uuid"
)

// TODO: probably put this in config/other?
var (
	validTypes []string = []string{"image/jpeg", "image/png", "image/gif"}
)

type ImageService struct {
	repo *ImageRepository
}

func NewService(r *ImageRepository) *ImageService {
	return &ImageService{repo: r}
}

// TODO: upon db failure, remove the saved image from disk
func (s *ImageService) UploadImage(ctx context.Context, r *http.Request, threadID, commentID int) (string, error) {
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
		return "", fmt.Errorf("image is not a valid type")
	}

	path := uuid.NewV7().String() + filepath.Ext(fileHeader.Filename)
	if err := createFile(file, path); err != nil {
		return "", fmt.Errorf("image creation failed: %w", err)
	}

	s.repo.Create(ctx, path, threadID, commentID)

	return path, nil
}

func validateImage(file multipart.File) (bool, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	filetype := http.DetectContentType(buffer)

	if !slices.Contains(validTypes, filetype) {
		return false, nil
	}

	return true, nil
}

func createFile(file multipart.File, filename string) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// TODO: remove hardcoded path
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", filename))
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
