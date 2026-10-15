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
// Returns an error combining all validation errors (if any).
func validateRegistrationCredentials(user *User, password string) error {
	var allErrs []error

	if !validateUsernameFormat(user.Username) {
		allErrs = append(allErrs, errs.ErrInvalidUsername)
	}

	if !validateEmailFormat(user.Email, emailRX) {
		allErrs = append(allErrs, errs.ErrInvalidEmail)
	}

	if !validatePasswordLength(password) {
		allErrs = append(allErrs, errs.ErrInvalidPassword)
	}

	return errors.Join(allErrs...)
}

// validateLoginCredentials validates the submitted email and password for correct format.
// Returns an error combining all validation errors (if any).
func validateLoginCredentials(input LoginSubmission) error {
	var allErrs []error

	if !validateEmailFormat(input.Email, emailRX) {
		allErrs = append(allErrs, errs.ErrInvalidEmail)
	}

	if !validatePasswordLength(input.Password) {
		allErrs = append(allErrs, errs.ErrInvalidPassword)
	}

	return errors.Join(allErrs...)
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
