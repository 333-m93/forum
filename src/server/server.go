package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"forum.com/m/src/database"
)

var dbConn *sql.DB

func StartServer() {

	// 🔥 DB CONNECTION (Render DATABASE_URL)
	var err error
	dbConn, err = database.Connect()
	if err != nil {
		log.Fatalf("❌ Erreur DB : %v", err)
	}

	// 🚀 PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	// 📡 ROUTES
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
		log.Fatalf("❌ Erreur serveur : %v", err)
	}
}

//
// =====================
// HANDLERS
// =====================
//

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
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Chat"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := "<div class='chat-header'><h2>" +
		template.HTMLEscapeString(name) +
		"</h2></div>"

	html += "<div class='chat-messages'>" +
		"<p>Bienvenue dans le chat <b>" +
		template.HTMLEscapeString(name) +
		"</b></p></div>"

	w.Write([]byte(html))
}
