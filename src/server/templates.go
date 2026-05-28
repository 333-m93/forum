package server

import "html/template"

var authTemplates = template.Must(template.ParseGlob("templates/*.html"))
