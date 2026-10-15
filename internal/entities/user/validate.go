package user

import (
	"errors"
	"forum/internal/errs"
	"regexp"
)

var (
	emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")
)

// validateRegistrationCredentials validates the submitted username, email, and password for user registration.
// Returns false if any of the validation checks fails.
func validateRegistrationCredentials(user *User, password string) error {
	var e []error

	if !validateUsernameFormat(user.Username) {
		e = append(e, errs.ErrInvalidUsername)
	}

	if !validatePasswordLength(password) {
		e = append(e, errs.ErrInvalidPassword)
	}

	if !validateEmailFormat(user.Email, emailRX) {
		e = append(e, errs.ErrInvalidEmail)
	}

	return errors.Join(e...)
}

// validateUsernameFormat checks that the username has 2-20 characters.
// Characters must be a letter (a-z, case insensitive), number, or underscore.
// Username must start with a letter. Returns false if validation fails.
func validateUsernameFormat(username string) bool {
	if len(username) < 2 || len(username) > 20 {
		return false
	}

	// First character must be a letter
	if !isLetter(rune(username[0])) {
		return false
	}

	for _, char := range username {
		if !(isLetter(char) || isNumber(char) || char == '_') {
			return false
		}
	}

	return true
}

// validatePasswordLength checks that the password has 8-128 characters. Returns false if validation fails.
func validatePasswordLength(password string) bool {
	n := len([]rune(password))

	if n < 8 || n > 128 {
		return false
	}

	return true
}

// validateEmailFormat checks that the email matches the email regexp. Returns false if validation fails.
func validateEmailFormat(email string, rx *regexp.Regexp) bool {
	return rx.MatchString(email)
}

func isLetter(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}

func isNumber(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}
