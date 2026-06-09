package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"forum.com/m/src/config"
	"forum.com/m/src/database"
)

var dbConn *sql.DB

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name       string
		Categories []string
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
	}

	if err := authTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Chat"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := "<div class=\"chat-header\"><h2 style='margin:0;color:#ffd8a8;'>" +
		template.HTMLEscapeString(name) + "</h2></div>"

	html += "<div class=\"chat-messages\" style='margin-top:14px;'>"
	html += "<div class=\"chat-message\">Bienvenue dans le chat <strong>" +
		template.HTMLEscapeString(name) + "</strong>.</div>"
	html += "<div class=\"chat-message\">Ici les messages apparaîtront.</div>"
	html += "</div>"

	html += "<div class=\"chat-footer\" style='margin-top:16px;'>" +
		"<form class=\"chat-form\" onsubmit=\"event.preventDefault();return false;\">" +
		"<input class=\"chat-input\" placeholder=\"Écrire un message...\" type=\"text\">" +
		"<button class=\"btn primary chat-send\">Envoyer</button></form></div>"

	w.Write([]byte(html))
}

func StartServer() {
	cfg := config.Load()

	// ✅ Render PORT priority
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	// 🔥 DB CONNECTION
	var err error
	dbConn, err = database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Erreur de connexion à la DB : %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/chat", chatHandler)
	mux.HandleFunc("/api/messages", messagesAPIHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	log.Printf("🚀 Serveur démarré sur %s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("❌ Erreur du serveur : %v", err)
	}
}
