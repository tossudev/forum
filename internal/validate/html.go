package validate

import (
	"forum/internal/errs"
	"html/template"
	"slices"
	"strings"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

var (
	HTMLPolicy          *bluemonday.Policy
	InvisibleCharacters = []rune{
		'\u2800', '\u3164',
	}
)

func CreatePolicy() {
	// Restrict all HTML initially
	HTMLPolicy = bluemonday.StrictPolicy()

	// Permits the "dir", "id", "lang", "title", "br" attributes
	HTMLPolicy.AllowStandardAttributes()
	HTMLPolicy.AllowElements("br")

	HTMLPolicy.AllowImages()
}

func ParseHTML(input string) (template.HTML, error) {
	input = strings.ReplaceAll(input, "\n", "<br>")

	out := HTMLPolicy.Sanitize(input)
	if out == "" {
		return template.HTML(""), errs.ErrHTMLSanitizationFail
	}

	return template.HTML(out), nil
}

// Remove non-printable (invisible) Unicode characters and leading & trailing whitespace
func TrimEmpty(input string) string {
	input = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) && !slices.Contains(InvisibleCharacters, r) {
			return r
		}
		return -1
	}, input)

	return strings.TrimSpace(input)
}
