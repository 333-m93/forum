package main

import (
	"log"
	"os"

	"forum.com/m/src/server"
)

func main() {

	log.Println("Starting forum server...")

	server.StartServer()

	_ = os.Getenv("PORT")
}
