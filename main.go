package main

import (
	"log"

	"forum.com/m/src/server"
)

func main() {
	log.Println("🚀 Starting forum server...")

	server.StartServer()
}
