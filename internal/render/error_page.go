package render

import (
	"errors"
	"forum/internal/errs"
	"log/slog"
	"net/http"
)

type ErrorPage struct {
	Username string
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
