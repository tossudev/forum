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
	// TODO: Validate user fields and password

	// Check if username already exists
	if exists, err := s.repo.UsernameExists(ctx, user.Username); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("%w: username already exists", errs.ErrDuplicate)
	}

	// Check if email already exists
	if exists, err := s.repo.EmailExists(ctx, user.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("%w: email already exists", errs.ErrDuplicate)
	}

	// Set password
	pw, err := password.New(inputPassword)
	if err != nil {
		return nil, err
	}

	user.Password = pw

	return s.repo.AddUser(ctx, user)
}

func (s *UserService) Authenticate(ctx context.Context, input CredentialsSubmission) (int, error) {
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
