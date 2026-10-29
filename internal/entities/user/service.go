package user

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"uuid"

	"forum/internal/config"
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/password"
	"forum/internal/validate"
)

type UserService struct {
	repo *UserRepo
}

func NewService(repo *UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, user *User, inputPassword string) (*User, error) {
	// Server-side validation of credentials in case a potentially malicious actor bypasses validation check in UI
	if err := validateRegistrationCredentials(user, inputPassword); err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	// Check if username already exists
	if exists, err := s.repo.UsernameExists(ctx, user.Username); err != nil {
		return nil, err
	} else if exists {
		return nil, errs.ErrDuplicateUsername
	}

	// Check if email already exists
	if exists, err := s.repo.EmailExists(ctx, user.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, errs.ErrDuplicateEmail
	}

	// Set password
	pw, err := password.New(inputPassword)
	if err != nil {
		return nil, err
	}

	user.Password = pw

	return s.repo.AddUser(ctx, user)
}

// Authenticate checks if the submitted user credentials for login is valid.
// Returns the user ID and nil error if valid. Failed authentication returns 0 and an error.
func (s *UserService) Authenticate(ctx context.Context, input LoginSubmission) (int, error) {
	// Validate email format and password length before further processing
	if err := validateLoginCredentials(input); err != nil {
		return 0, fmt.Errorf("%w: %w", errs.ErrInvalidUserInput, err)
	}

	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return 0, err
	}

	// Check if password hash matches
	matches, err := password.VerifyPassword(input.Password, user.Password.Hash, user.Password.Salt)
	if err != nil {
		return 0, err
	}

	if !matches {
		return 0, fmt.Errorf("%w: invalid credentials", errs.ErrUnauthorized)
	}

	return user.ID, nil
}

func (s *UserService) GetUserProfileByUsername(ctx context.Context, username string) (*UserProfile, error) {
	return s.repo.GetUserProfileByUsername(ctx, username)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetThreadsByUsername(ctx context.Context, username string, pagination pagination.Pagination) ([]ProfileThread, int, error) {
	threads, totalCount, err := s.repo.GetThreadsByUsername(ctx, username, pagination)
	if err != nil {
		return nil, 0, err
	}

	for i, t := range threads {
		threads[i].TitleHTML, err = validate.ParseHTML(t.Title)
		if err != nil {
			slog.Error("thread title validtion failed", "err:", err, "thread_id:", t.ID, "title:", t.Title)
		}

		threads[i].BodyHTML, err = validate.ParseHTML(t.Body)
		if err != nil {
			slog.Error("thread body validation failed", "err:", err, "thread_id:", t.ID, "body:", t.Body)
		}
	}

	return threads, totalCount, nil
}

func (s *UserService) GetCommentsByUsername(ctx context.Context, username string, pagination pagination.Pagination) ([]ProfileComment, int, error) {
	comments, totalCount, err := s.repo.GetCommentsByUsername(ctx, username, pagination)
	if err != nil {
		return nil, 0, err
	}

	for i, c := range comments {
		comments[i].BodyHTML, err = validate.ParseHTML(c.Body)
		if err != nil {
			slog.Error("comment validation failed", "err:", err, "comment_id:", c.ID, "body:", c.Body)
		}
	}

	return comments, totalCount, nil
}

func (s *UserService) GetLikedThreadsByUsername(ctx context.Context, username string, pagination pagination.Pagination) ([]ProfileThread, int, error) {
	threads, totalCount, err := s.repo.GetLikedThreadsByUsername(ctx, username, pagination)
	if err != nil {
		return nil, 0, err
	}

	for i, t := range threads {
		threads[i].TitleHTML, err = validate.ParseHTML(t.Title)
		if err != nil {
			slog.Error("thread title validtion failed", "err:", err, "thread_id:", t.ID, "title:", t.Title)
		}

		threads[i].BodyHTML, err = validate.ParseHTML(t.Body)
		if err != nil {
			slog.Error("thread body validation failed", "err:", err, "thread_id:", t.ID, "body:", t.Body)
		}
	}

	return threads, totalCount, nil
}

func (s *UserService) GetDislikedThreadsByUsername(ctx context.Context, username string, pagination pagination.Pagination) ([]ProfileThread, int, error) {
	threads, totalCount, err := s.repo.GetDislikedThreadsByUsername(ctx, username, pagination)
	if err != nil {
		return nil, 0, err
	}

	for i, t := range threads {
		threads[i].TitleHTML, err = validate.ParseHTML(t.Title)
		if err != nil {
			slog.Error("thread title validtion failed", "err:", err, "thread_id:", t.ID, "title:", t.Title)
		}

		threads[i].BodyHTML, err = validate.ParseHTML(t.Body)
		if err != nil {
			slog.Error("thread body validation failed", "err:", err, "thread_id:", t.ID, "body:", t.Body)
		}
	}

	return threads, totalCount, nil
}

func (s *UserService) ChangeAvatar(ctx context.Context, r *http.Request, userID int) error {
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return fmt.Errorf("change avatar formfile err: %w", err)
	}
	defer file.Close()

	valid, err := validateFiletype(file)
	if err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}

	if !valid {
		return errs.ErrInvalidFiletype
	}

	imageName := uuid.NewV7().String() + filepath.Ext(fileHeader.Filename)
	imagePath := filepath.Join(config.UploadsPath, imageName)
	if err := createFile(file, imagePath); err != nil {
		return fmt.Errorf("image creation failed: %w", err)
	}

	return s.repo.ChangeAvatar(ctx, imageName, userID)
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
