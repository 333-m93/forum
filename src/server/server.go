package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"forum.com/m/src/config"
	"forum.com/m/src/database"
)

var dbConn *sql.DB

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name       string
		Categories []string
	}{
		Name:       "Le Dojo",
		Categories: []string{"Chat général", "MMA", "Boxe Anglaise", "Muay Thai", "Jujitsu Brésilien", "Grappling", "Autres sports de combat"},
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
	// Simple HTML fragment for the chat pane. In future this can be rendered from templates.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := "<div class=\"chat-header\"><h2 style='margin:0;color:#ffd8a8;'>" + template.HTMLEscapeString(name) + "</h2></div>"
	html += "<div class=\"chat-messages\" style='margin-top:14px;'>"
	// placeholder messages
	html += "<div class=\"chat-message\">Bienvenue dans le chat <strong>" + template.HTMLEscapeString(name) + "</strong>.</div>"
	html += "<div class=\"chat-message\">Ici les messages apparaîtront.</div>"
	html += "</div>"
	html += "<div style='margin-top:16px;'><form onsubmit=\"event.preventDefault();return false;\"><input placeholder=\"Écrire un message...\" style=\"width:76%;padding:10px;border-radius:10px;border:1px solid rgba(255,255,255,0.06);background:#0b0b0d;color:#fff;\"><button class=\"btn primary\" style=\"margin-left:8px;\">Envoyer</button></form></div>"
	w.Write([]byte(html))
}

func StartServer() {
	cfg := config.Load()
	var err error
	dbConn, err = database.Connect(cfg)
	if err != nil {
		log.Fatalf("Erreur de connexion à la DB : %v", err)
	}
	defer dbConn.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/chat", chatHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Serveur démarré sur %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}
