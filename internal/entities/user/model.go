package user

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	RoleID    int       `json:"role_id"`
}

type Password struct {
	plaintext *string
	hash      []byte
	salt      []byte
}
