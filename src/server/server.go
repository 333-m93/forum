package server

import (
	"database/sql"
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
		Name:       "CombatArena",
		Categories: []string{"Chat général", "MMA", "Boxe Anglaise", "Muay Thai", "Jujitsu Bresilien", "Grappling", "Autres arts martiaux"},
	}
	if err := authTemplates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
	}
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
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Serveur démarré sur %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}
