package thread

import (
	"errors"
	"forum/internal/errs"
	"log/slog"
	"net/http"
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

func getThreadCreationFormFieldErrors(err error) (map[string]string, int) {
	fieldErrors := make(map[string]string)

	if err == nil {
		return fieldErrors, http.StatusOK
	}

	switch {
	case errors.Is(err, errs.ErrInvalidUserInput):
		if errors.Is(err, errs.ErrInvalidThreadTitle) {
			fieldErrors["TitleErr"] = errs.ErrInvalidThreadTitle.Error()
		}
		if errors.Is(err, errs.ErrInvalidThreadBody) {
			fieldErrors["BodyErr"] = errs.ErrInvalidThreadBody.Error()
		}
		if errors.Is(err, errs.ErrHTMLSanitizationFail) {
			fieldErrors["OtherErr"] = errs.ErrHTMLSanitizationFail.Error()
		}
		return fieldErrors, http.StatusBadRequest
	case errors.Is(err, errs.ErrUnauthorized):
		fieldErrors["OtherErr"] = errs.ErrInvalidCredentials.Error()
		return fieldErrors, http.StatusUnauthorized
	default:
		fieldErrors["OtherErr"] = errs.ErrInternal.Error()
		slog.Error("internal server error:", "err", err)
		return fieldErrors, http.StatusInternalServerError
	}
}
