// cmd/routes.go
package main

import (
	"net/http"
	"storytelling-backend/internal/api"

	"github.com/gorilla/mux"
)

// Route defines the structure for a route in the application.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// SetupRoutes initializes the routes for the application.
func SetupRoutes(router *mux.Router) {
	routes := []Route{
		{"POST", "/create-room", api.CreateRoomHandler},
		{"POST", "/join-room", api.JoinRoomHandler},
		{"POST", "/start-game/{room_id}", api.StartGameHandler},
		{"POST", "/submit-line", api.SubmitLineHandler},
		{"GET", "/get-story", api.GetStoryHandler},
		{"GET", "/ws", api.WebSocketHandler},
	}

	for _, route := range routes {
		router.Handle(route.Path, corsMiddleware(route.Handler)).Methods(route.Method)
	}
}

// corsMiddleware is used to handle CORS for the application.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins (or specify your front-end URL)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // Allow specific methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allow specific headers

		if r.Method == http.MethodOptions { // Handle preflight request
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
