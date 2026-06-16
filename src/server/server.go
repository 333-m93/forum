package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"forum.com/m/src/database"
)

var dbConn *sql.DB

func StartServer() {

	var err error
	dbConn, err = database.Connect()
	if err != nil {
		log.Fatalf("❌ DB error: %v", err)
	}

	if err := database.Migrate(dbConn); err != nil {
		log.Fatalf("❌ Migration error: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	mux := http.NewServeMux()

	// UI
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/chat", chatHandler)

	// API
	mux.HandleFunc("/api/messages", messagesAPIHandler)

	// AUTH
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	// STATIC
	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	log.Printf("🚀 Server running on %s", port)

	log.Fatal(http.ListenAndServe(port, mux))
}
