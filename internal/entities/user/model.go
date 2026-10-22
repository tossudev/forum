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

// LoginSubmission contains the user input during login before credential verification
type LoginSubmission struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProfilePage struct {
	Username    string
	Role        string
	DateCreated string
}

type RegistrationPage struct {
	Username      string // Logged in user, if any (for toggling login/logout button)
	InputUsername string // From registration form
	InputEmail    string
	UsernameErr   string
	EmailErr      string
	PasswordErr   string
	OtherErr      string
}

type LoginPage struct {
	Username    string // Only needed because this field is required by the base layout
	Email       string
	EmailErr    string
	PasswordErr string
	OtherErr    string
}

// ProfileThread contains thread data for displaying user's posted threads.
// Contains the same properties as thread.Thread.
// TODO: Perhaps move thread.Thread into a shared models package so the same struct type can be reused without circular imports.
type ProfileThread struct {
	ID          int
	Title       string
	Body        string
	DateCreated string
	AuthorID    int
	CategoryID  int
	ReplyCount  int
}

type ProfileComment struct {
	ID          int
	Body        string
	DateCreated string
	ThreadID    int
	AuthorID    int
}
