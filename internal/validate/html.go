package validate

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	HTMLPolicy *bluemonday.Policy = bluemonday.UGCPolicy()
)

func ParseHTML(input string) (template.HTML, error) {
	input = strings.ReplaceAll(input, "\n", "<br>")
	input = strings.ReplaceAll(input, "\t", "&emsp;")

	out := HTMLPolicy.Sanitize(input)
	if out == "" {
		return template.HTML(""), fmt.Errorf("HTML sanitization failed")
	}

	return template.HTML(out), nil
}
