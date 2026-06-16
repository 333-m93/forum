package server

import (
	"html/template"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := GetSessionUser(r, dbConn)

	data := struct {
		Name       string
		Categories []string
		User       *User
	}{
		Name: "Le Dojo",
		Categories: []string{
			"Chat général",
			"MMA",
			"Boxe Anglaise",
			"Muay Thai",
			"Jujitsu Brésilien",
			"Grappling",
			"Autres sports de combat",
		},
		User: user,
	}

	err := authTemplates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), 500)
	}
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
