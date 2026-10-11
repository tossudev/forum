package user

import (
	"regexp"
)

var (
	emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// validateRegistrationCredentials validates the submitted username, email, and password for user registration.
// Returns false if any of the validation checks fails.
func validateRegistrationCredentials(user *User, password string) bool {
	if !validateUsername(user.Username) {
		return false
	}

	if !validatePassword(password) {
		return false
	}

	if !validateEmail(user.Email, emailRX) {
		return false
	}

	return true
}

// validateUsername checks that the username has 2-20 characters.
// Characters must be a letter (a-z, case insensitive), number, or underscore.
// Username must start with a letter. Returns false if validation fails.
func validateUsername(username string) bool {
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

// validatePassword checks that the password has 8-128 characters. Returns false if validation fails.
func validatePassword(password string) bool {
	n := len([]rune(password))

	if n < 8 || n > 128 {
		return false
	}

	return true
}

// validateEmail checks that the email matches the email regexp. Returns false if validation fails.
func validateEmail(email string, rx *regexp.Regexp) bool {
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
