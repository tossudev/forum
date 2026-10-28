package user

import (
	"forum/internal/password"
	"html/template"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  password.Password
	CreatedAt string
	RoleID    int
}

// Data needed for user profile pages
type UserProfile struct {
	Username  string // Username of the profile being requested
	Avatar    string
	CreatedAt string
}

// LoginSubmission contains the user input during login before credential verification
type LoginSubmission struct {
	Email    string
	Password string
}

type ProfilePage struct {
	ProfileUsername string
	Username        string // logged in user
	Role            string
	DateCreated     string
}

type RegistrationPage struct {
	Username      string // For base layout
	InputUsername string
	InputEmail    string
	FieldErrors   map[string]string
}

type LoginPage struct {
	Username    string // For base layout
	Email       string
	FieldErrors map[string]string
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
	TitleHTML   template.HTML
	BodyHTML    template.HTML
}

type ProfileComment struct {
	ID          int
	Body        string
	DateCreated string
	ThreadID    int
	AuthorID    int
	BodyHTML    template.HTML
}
