package user

import (
	"errors"
	"forum/internal/errs"
	"log/slog"
	"net/http"
	"regexp"
)

var (
	emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")
)

// getRegistrationFieldErrors returns a map of field errors and corresponding error messages for display.
// Also returns the http status code.
func getRegistrationFieldErrors(err error) (map[string]string, int) {
	fieldErrors := make(map[string]string)

	if err == nil {
		return fieldErrors, http.StatusOK
	}

	// Username field errors
	switch {
	case errors.Is(err, errs.ErrInvalidUsername):
		fieldErrors["UsernameErr"] = "Username must be 2-20 characters, start with a letter, and only contain letters, numbers, or underscores (_)."
	case errors.Is(err, errs.ErrDuplicateUsername):
		fieldErrors["UsernameErr"] = "Username already taken."
	}

	// Email field errors
	switch {
	case errors.Is(err, errs.ErrInvalidEmail):
		fieldErrors["EmailErr"] = "Enter a valid email address, such as example@email.com."
	case errors.Is(err, errs.ErrDuplicateEmail):
		fieldErrors["EmailErr"] = "Email already taken."
	}

	// Password field errors
	if errors.Is(err, errs.ErrInvalidPassword) {
		fieldErrors["PasswordErr"] = "Password must be 8-128 characters."
	}

	if len(fieldErrors) > 0 {
		return fieldErrors, http.StatusBadRequest
	}

	// None of the above errors
	fieldErrors["OtherErr"] = "Sorry, we're unable to process user registration at the moment. Please try again later."
	slog.Error("user registration error:", "err", err)

	return fieldErrors, http.StatusInternalServerError
}

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
	return len(email) >= 6 && len(email) <= 254 && rx.MatchString(email)
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
