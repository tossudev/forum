package user

import (
	"forum/internal/password"
	"time"
)

type User struct {
	ID        int               `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Password  password.Password `json:"-"`
	CreatedAt time.Time         `json:"created_at"`
	RoleID    int               `json:"role_id"`
}
