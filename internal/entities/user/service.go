package user

import (
	"context"
	"fmt"
	"forum/internal/errs"
	"forum/internal/password"
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

func (s *UserService) GetUserByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}
