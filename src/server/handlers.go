package server

import (
	"html/template"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello forum"))
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Chat"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := "<h1>" + template.HTMLEscapeString(name) + "</h1>"
	w.Write([]byte(html))
}
