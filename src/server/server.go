package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bonjour depuis le serveur Go !")
}

func main() {
	http.HandleFunc("/", helloHandler)

	addr := ":8080"
	log.Printf("Serveur démarré sur %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}
