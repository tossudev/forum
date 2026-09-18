package utils

import (
	"html/template"
)

var Templates = template.Must(
	template.ParseGlob(
		"web/tmpl/*html",
	),
)
