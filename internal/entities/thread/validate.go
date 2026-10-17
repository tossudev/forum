package thread

import (
	"errors"
	"forum/internal/errs"
	"unicode/utf8"
)

func validateThreadTitle(title string) bool {
	numChars := utf8.RuneCountInString(title)

	if numChars < 1 || numChars > 100 {
		return false
	}

	return true
}

func validateThreadBody(body string) bool {
	numChars := utf8.RuneCountInString(body)

	if numChars < 1 || numChars > 5000 {
		return false
	}

	return true
}

func validateThreadContent(title, body string) error {
	var allErrs []error

	if !validateThreadTitle(title) {
		allErrs = append(allErrs, errs.ErrInvalidThreadTitle)
	}

	if !validateThreadBody(body) {
		allErrs = append(allErrs, errs.ErrInvalidThreadBody)
	}

	return errors.Join(allErrs...)
}
