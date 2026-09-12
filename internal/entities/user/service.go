package user

import (
	"context"
	"forum/internal/password"
)

type UserService struct {
	repo *UserRepo
}

func NewService(repo *UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, user *User, inputPassword string) error {
	// TODO: Validate user fields

	// Set password
	pw, err := password.New(inputPassword)
	if err != nil {
		return err
	}

	user.Password = pw

	return nil
}
