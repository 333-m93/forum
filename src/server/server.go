package server

import (
	"fmt"
	"log"
	"net/http"

	"forum.com/m/src/config"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bonjour depuis le serveur Go !")
}

func StartServer() {
	cfg := config.Load()

	http.HandleFunc("/", helloHandler)

	log.Printf("Serveur démarré sur %s\n", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, nil); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}
