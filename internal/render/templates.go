package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"forum/internal/errs"
)

type Renderer struct {
	Templates map[string]*template.Template
}

func NewRenderer() *Renderer {
	return &Renderer{
		Templates: parseTemplates(),
	}
}

func (r *Renderer) RenderPage(w http.ResponseWriter, templateName string, data any) {
	tmpl, ok := r.Templates[templateName]
	if !ok {
		errs.WriteError(w, fmt.Errorf("Couldn't find template: %s", templateName))
		return
	}

	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base.html", data) // use buffer to prevent partial writing of response in case of execution failure
	if err != nil {
		errs.WriteError(w, err)
		return
	}

	buf.WriteTo(w)
}

// Prepare all templates and cache page templates in map
// TODO: probably want to change hardcoded paths?
func parseTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)

	// Create base template (base + components)
	baseTmpl, err := template.ParseFiles("web/tmpl/base.html")
	if err != nil {
		slog.Error(fmt.Sprintf("Template parsing failure: %v", err))
	}

	// Get all files from pages folder -> []DirEntry
	files, err := os.ReadDir("web/tmpl/pages")
	if err != nil {
		slog.Error(fmt.Sprintf("Read directory failure: %v", err))
	}

	// Create page templates: For each page file, parse with base template and cache in map
	for _, file := range files {
		tmpl, err := baseTmpl.Clone() // clone to avoid mutation of base
		if err != nil {
			slog.Error(fmt.Sprintf("Template cloning failure: %v", err))
		}

		pageTmpl, err := tmpl.ParseFiles("web/tmpl/pages/" + file.Name()) // file.Name() -> e.g. "home.html"
		if err != nil {
			slog.Error(fmt.Sprintf("Template parsing failure: %v", err))
		}
		templates[file.Name()] = pageTmpl
	}

	return templates
}
