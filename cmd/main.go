// cmd/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"storytelling-backend/config"
	"storytelling-backend/internal/game"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Create a new router
	router := mux.NewRouter()

	// Setup routes
	SetupRoutes(router)
	game.RoomManagerInstance = game.NewRoomManager()
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
