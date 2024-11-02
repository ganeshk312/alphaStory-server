// internal/models/player_connection.go
package models

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// PlayerConnection represents a player with an associated WebSocket connection in a room.
type PlayerConnection struct {
	Conn       *websocket.Conn
	PlayerName string
	RoomID     string
}

// NewPlayerConnection initializes a new player connection.
func NewPlayerConnection(conn *websocket.Conn, roomID, playerName string) *PlayerConnection {
	return &PlayerConnection{
		Conn:       conn,
		RoomID:     roomID,
		PlayerName: playerName,
	}
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// // Listen listens for incoming messages from the player and manages disconnections.
// func (p *PlayerConnection) Listen(handleMessage func(Message)) {
// 	defer func() {
// 		p.Conn.Close()
// 	}()

// 	for {
// 		var msg Message
// 		err := p.Conn.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("Player %s disconnected: %v", p.PlayerName, err)
// 			break
// 		}

// 		// Pass the message to the provided message handler.
// 		handleMessage(msg)
// 	}
// }

// Listen listens for incoming messages from the player and manages disconnections.
func (p *PlayerConnection) Listen(room *Room) {
	defer func() {
		// room, err := game.RoomManagerInstance.GetRoom(p.RoomID)
		// if err == nil {
		room.RemovePlayer(p.PlayerName)
		room.BroadcastMessage(p.PlayerName + " has left the game.")
		if room.Status != "completed" {
			room.BroadcastTurn() // Notify the next player if a player disconnects during a live game.
		}
		if len(room.Players) == 0 {
			room.Status = "completed" //Handle cleanup if all players leave
		}
		// }
		p.Conn.Close()
	}()

	for {
		var msg Message
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Player %s disconnected: %v", p.PlayerName, err)
			break
		}

		// Process different types of incoming messages
		switch msg.Type {
		case "SUBMIT_LINE":
			// room, err := game.RoomManagerInstance.GetRoom(p.RoomID)
			room.HandleSubmitLine(p.PlayerName, msg.Content)

		case "START_GAME":
			// room, err := game.RoomManagerInstance.GetRoom(p.RoomID)
			if p.PlayerName == room.Host {
				room.StartGame()
			} else {
				log.Printf("Player %s is not the host and cannot start the game", p.PlayerName)
				p.SendMessage(Message{Type: "ERROR", Content: "Only the host can start the game"})
			}
		default:
			log.Printf("Unhandled message type from %s: %s", p.PlayerName, msg.Type)
		}
	}
}

func (pc *PlayerConnection) Send(msg Message) error {
	return pc.Conn.WriteJSON(msg)
}

// SendStoryUpdate sends the current story to the player.
func (pc *PlayerConnection) SendStoryUpdate(story []string) {
	update := map[string]interface{}{
		"type":  "story_update",
		"story": story,
	}
	pc.SendMessage(update)
}

// SendMessage sends a message to the player over WebSocket.
func (pc *PlayerConnection) SendMessage(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message for player %s: %v", pc.PlayerName, err)
		return
	}
	err = pc.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("Failed to send message to player %s: %v", pc.PlayerName, err)
	}
}
