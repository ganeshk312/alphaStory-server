// internal/api/http_handler.go
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"storytelling-backend/internal/game"
	"storytelling-backend/pkg/utils"

	"github.com/gorilla/mux"
)

// Request payload structures
type CreateRoomRequest struct {
	StoryName  string `json:"story_name"`
	PlayerName string `json:"player_name"`
}

type JoinRoomRequest struct {
	RoomID     string `json:"room_id"`
	PlayerName string `json:"player_name"`
}

type AddLineRequest struct {
	RoomID     string `json:"room_id"`
	PlayerName string `json:"player_name"`
	Line       string `json:"line"`
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateRoomHandler called")
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("Creating room with story name %s by player %s", req.StoryName, req.PlayerName)

	roomID := utils.GenerateRoomID() // Function to generate a unique room ID
	room, err := game.RoomManagerInstance.CreateRoom(roomID, req.PlayerName)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the host as a player
	room.AddPlayer(req.PlayerName)
	log.Printf("Room %s created successfully with host player %s", roomID, req.PlayerName)

	// Return the room ID in the response
	response := map[string]string{"room_id": room.ID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("JoinRoomHandler called")
	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Printf("Joining room %s with player %s", req.RoomID, req.PlayerName)

	room, err := game.RoomManagerInstance.AddPlayerToRoom(req.RoomID, req.PlayerName)
	if err != nil {
		log.Printf("Error joining room: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(room); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("Player %s joined room %s successfully", req.PlayerName, req.RoomID)
}

// // AddLineToStoryHandler allows a player to add a line to the story in a specified room
// func AddLineToStoryHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("AddLineToStoryHandler called")
// 	var req AddLineRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		log.Printf("Error decoding request: %v", err)
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}
// 	log.Printf("Adding line to room %s by player %s", req.RoomID, req.PlayerName)

// 	room, err := game.RoomManagerInstance.AddLineToStory(req.RoomID, req.PlayerName, req.Line)
// 	if err != nil {
// 		log.Printf("Error adding line to story: %v", err)
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(room); err != nil {
// 		log.Printf("Error encoding response: %v", err)
// 	}
// 	log.Printf("Line added to room %s by player %s successfully", req.RoomID, req.PlayerName)
// }

func GetStoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetStoryHandler called")
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		log.Println("Room ID is required")
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	log.Printf("Retrieving story for room %s", roomID)

	story, err := game.RoomManagerInstance.GetStory(roomID)
	if err != nil {
		log.Printf("Error retrieving story: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(story); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Printf("Story for room %s retrieved successfully", roomID)
}

// Obsoleted by WebSocketHandler
func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("StartGameHandler called")
	roomID := mux.Vars(r)["room_id"]
	log.Printf("Starting game for room %s", roomID)

	room, err := game.RoomManagerInstance.GetRoom(roomID)
	if err != nil {
		log.Printf("Error finding room: %v", err)
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	room.StartGame()
	log.Printf("Game started for room %s", roomID)
	w.WriteHeader(http.StatusOK)
}

// Obsoleted by WebSocketHandler
func SubmitLineHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SubmitLineHandler called")
	var req struct {
		RoomID     string `json:"room_id"`
		PlayerName string `json:"player_name"`
		Line       string `json:"line"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Printf("Submitting line to room %s by player %s", req.RoomID, req.PlayerName)

	room, err := game.RoomManagerInstance.GetRoom(req.RoomID)
	if err != nil {
		log.Printf("Error finding room: %v", err)
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if err := room.AddLine(req.PlayerName, req.Line); err != nil {
		log.Printf("Error adding line to story: %v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Line submitted to room %s by player %s successfully", req.RoomID, req.PlayerName)
}
