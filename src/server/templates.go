package server

import (
	"html/template"
	"log"
	"strings"
)

var authTemplates *template.Template

func init() {
	var err error

	funcMap := template.FuncMap{
		"uppercase": strings.ToUpper,
	}

	authTemplates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Erreur chargement templates: %v", err)
	}
}
