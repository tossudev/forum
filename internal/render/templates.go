package render

import (
	"bytes"
	"errors"
	"fmt"
	"forum/internal/errs"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

const (
	templatesDir string = "web/tmpl"
)

type Renderer struct {
	Templates map[string]*template.Template
}

func NewRenderer() *Renderer {
	return &Renderer{
		Templates: parseTemplates(),
	}
}

type ErrorPage struct {
	Username string
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
		errs.WriteError(w, fmt.Errorf("template execution failed: %w", err))
		return
	}

	buf.WriteTo(w)
}

// WriteErrorPage renders a generic error page based on the error along with writing the http status code
func (r *Renderer) WriteErrorPage(w http.ResponseWriter, err error, pageData ErrorPage) {
	switch {
	case errors.Is(err, errs.ErrInvalidUserInput):
		w.WriteHeader(http.StatusBadRequest)
		r.RenderPage(w, "error-400.html", pageData)
	case errors.Is(err, errs.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		r.RenderPage(w, "error-404.html", pageData)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		r.RenderPage(w, "error-500.html", pageData)
		slog.Error("internal server error", "err", err)
	}
}

// Prepare all templates and cache page templates in map
func parseTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)

	// Create base template (base + components)
	baseTmpl, err := template.ParseFiles(templatesDir + "/base.html")
	if err != nil {
		slog.Error(fmt.Sprintf("Template parsing failure: %v", err))
	}

	baseTmpl, err = baseTmpl.ParseGlob(templatesDir + "/partials/*.html")
	if err != nil {
		slog.Error(fmt.Sprintf("Template parsing failure: %v", err))
	}

	// Get all files from pages folder -> []DirEntry
	files, err := os.ReadDir(templatesDir + "/pages")
	if err != nil {
		slog.Error(fmt.Sprintf("Read directory failure: %v", err))
	}

	// Create page templates: For each page file, parse with base template and cache in map
	for _, file := range files {
		tmpl, err := baseTmpl.Clone() // clone to avoid mutation of base
		if err != nil {
			slog.Error(fmt.Sprintf("Template cloning failure: %v", err))
		}

		pageTmpl, err := tmpl.ParseFiles(templatesDir + "/pages/" + file.Name()) // file.Name() -> e.g. "home.html"
		if err != nil {
			slog.Error(fmt.Sprintf("Template parsing failure: %v", err))
		}
		templates[file.Name()] = pageTmpl
	}

	return templates
}
