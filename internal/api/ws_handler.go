// internal/api/ws_handler.go
package api

import (
	"net/http"
	"storytelling-backend/internal/game"
	"storytelling-backend/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; adjust in production for security
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Extract room ID and player name from query parameters
	roomID := r.URL.Query().Get("room_id")
	playerName := r.URL.Query().Get("player_name") // Assuming you're passing the player's name as well
	if roomID == "" || playerName == "" {
		http.Error(w, "Room ID and Player Name are required", http.StatusBadRequest)
		return
	}

	// Register the WebSocket connection with the room
	playerConn := models.NewPlayerConnection(conn, roomID, playerName) // Include playerName
	if err := game.RoomManagerInstance.AddConnectionToRoom(roomID, playerConn); err != nil {
		http.Error(w, "Failed to register connection:"+err.Error(), http.StatusInternalServerError)
		return
	}
	room, _ := game.RoomManagerInstance.GetRoom(roomID)
	room.BroadcastMessage(playerName + " joined the room.")
	// Handle incoming messages and player disconnects
	if len(room.Players) == room.TotalPlayers {
		room.BroadcastTurn() // Start the game when all players join
	}

	playerConn.Listen(room)
}
