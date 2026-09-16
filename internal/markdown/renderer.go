package main

import (
	"bytes"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type Renderer struct {
	parser parser.Parser
	policy *bluemonday.Policy
}

func NewRenderer() *Renderer {
	parser := parser.New()
	policy := bluemonday.UGCPolicy()

	return &Renderer{
		parser: parser,
		policy: policy,
	}
}

func (r *Renderer) RenderHTML(input string) (string, error) {
	html, err := r.ParseMarkdown(input)
	if err != nil {
		return "", err
	}

	rendered := r.SanitizeHTML(html)

	if rendered == "" {
		return "", fmt.Errorf("HTML sanitization failed")
	}

	return rendered, nil
}

func (r *Renderer) ParseMarkdown(input string) (string, error) {
	source := []byte(input)

	var buf bytes.Buffer
	p := parser.New()
	re := html.New()

	doc := p.Parse(source)
	if err := re.Render(&buf, source, doc); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *Renderer) SanitizeHTML(input string) string {
	return r.policy.Sanitize(input)
}
