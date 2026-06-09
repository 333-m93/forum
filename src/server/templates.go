package server

import (
	"html/template"
	"log"
)

var authTemplates *template.Template

func init() {
	var err error

	authTemplates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("❌ Erreur chargement templates: %v", err)
	}
}
