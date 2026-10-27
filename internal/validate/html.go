package validate

import (
	"forum/internal/errs"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	HTMLPolicy *bluemonday.Policy
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
