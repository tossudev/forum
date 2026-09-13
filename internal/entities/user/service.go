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

func (s *UserService) RegisterUser(ctx context.Context, user *User, inputPassword string) error {
	// TODO: Validate user fields and password

	// Check if username already exists
	if exists, err := s.repo.UsernameExists(ctx, user.Username); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("%w: username already exists", errs.ErrDuplicate)
	}

	// Check if email already exists
	if exists, err := s.repo.EmailExists(ctx, user.Email); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("%w: email already exists", errs.ErrDuplicate)
	}

	// Set password
	pw, err := password.New(inputPassword)
	if err != nil {
		return err
	}

	user.Password = pw

	return s.repo.AddUser(ctx, user)
}
