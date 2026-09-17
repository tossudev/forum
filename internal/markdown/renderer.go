package markdown

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type Renderer struct {
	policy *bluemonday.Policy
}

var renderer *Renderer = NewRenderer()

// Sanitizer uses default UGCPolicy
// https://pkg.go.dev/github.com/microcosm-cc/bluemonday?utm_source=godoc#UGCPolicy
func NewRenderer() *Renderer {
	return &Renderer{
		policy: bluemonday.UGCPolicy(),
	}
}

func RenderHTML(input string) (string, error) {
	return renderer.SanitizeHTML(string(blackfriday.Run([]byte(input))))
}

func (r *Renderer) SanitizeHTML(input string) (string, error) {
	out := r.policy.Sanitize(input)
	if out == "" {
		return "", fmt.Errorf("HTML sanitization failed")
	}

	return out, nil
}
