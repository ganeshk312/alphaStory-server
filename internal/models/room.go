// internal/models/room.go
package models

import (
	"errors"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Room represents a storytelling room with a unique ID, list of players, and the story.
type Room struct {
	ID           string
	Host         string
	Players      map[string]*PlayerConnection
	Story        []string
	TurnOrder    []string
	CurrentTurn  int
	Status       string
	Mutex        sync.Mutex
	TotalPlayers int
}

// NewRoom creates a new Room with a specified ID.
func NewRoom(roomID, host string) *Room {
	return &Room{
		ID:          roomID,
		Host:        host,
		Players:     make(map[string]*PlayerConnection),
		Story:       []string{},
		TurnOrder:   []string{},
		CurrentTurn: 0,
		Status:      "waiting",
	}
}

// AddPlayer adds a player connection to the room.
func (r *Room) AddPlayer(playerName string) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if _, exists := r.Players[playerName]; !exists {
		r.Players[playerName] = &PlayerConnection{PlayerName: playerName, RoomID: r.ID}
		r.TurnOrder = append(r.TurnOrder, playerName)
		return nil
	}
	return errors.New("player already exists")
}

func (r *Room) RemovePlayer(playerName string) {
	// r.Mutex.Lock()
	// defer r.Mutex.Unlock()

	delete(r.Players, playerName)
	for i, name := range r.TurnOrder {
		if name == playerName {
			r.TurnOrder = append(r.TurnOrder[:i], r.TurnOrder[i+1:]...)
			break
		}
	}
}

// AddConnection assigns a WebSocket connection to a player.
func (r *Room) AddConnection(conn *PlayerConnection) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	if r.Players[conn.PlayerName] == nil {
		return errors.New("player not found in room")
	}
	if r.Players[conn.PlayerName].Conn != nil {
		r.Players[conn.PlayerName].Conn.Close()
	}
	r.Players[conn.PlayerName] = conn
	return nil
}

// Start the game by setting the first player's turn
func (r *Room) StartGame() {
	// r.Mutex.Lock()
	// defer r.Mutex.Unlock()

	r.Status = "in_progress"
	r.CurrentTurn = 0
	r.BroadcastMessage("Game started! It's " + r.TurnOrder[r.CurrentTurn] + "'s turn.")
}

// Add a line to the story and move to the next turn
func (r *Room) AddLine(playerName, line string) error {
	// r.Mutex.Lock()
	// defer r.Mutex.Unlock()

	// Validate if it's the player's turn
	if r.TurnOrder[r.CurrentTurn] != playerName {
		return errors.New("it's not your turn")
	}

	if _, exists := r.Players[playerName]; !exists {
		return errors.New("player not found in room")
	}
	// Update story
	r.Story = append(r.Story, line)
	r.advanceTurn()

	return nil
}

// Move to the next player's turn, or end the game if all turns are completed
func (r *Room) advanceTurn() {
	r.CurrentTurn++
	if r.CurrentTurn >= len(r.TurnOrder) {
		r.Status = "completed"
		r.BroadcastMessage("Game completed! Final story: " + r.GetStory())
	} else {
		r.BroadcastMessage("It's " + r.TurnOrder[r.CurrentTurn] + "'s turn.")
	}
}

func (r *Room) BroadcastMessage(message string) {
	// r.Mutex.Lock()
	// defer r.Mutex.Unlock()

	for _, player := range r.Players {
		if player.Conn != nil {
			player.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

// GetStory returns the full story as a single string
func (r *Room) GetStory() string {
	return "Story: " + strings.Join(r.Story, " ")
}

// BroadcastStoryUpdate sends the updated story to all players in the room.
// func (r *Room) BroadcastStoryUpdate() {
// r.connMutex.Lock()
// defer r.connMutex.Unlock()

// 	for _, playerConn := range r.Players {
// 		if playerConn != nil && playerConn.Conn != nil {
// 			playerConn.SendStoryUpdate(r.Story)
// 		}
// 	}
// }

func (r *Room) BroadcastTurn() {
	// r.Mutex.Lock()
	// defer r.Mutex.Unlock()

	currentPlayer := r.TurnOrder[r.CurrentTurn]
	message := Message{Type: "TURN", Content: currentPlayer + "'s turn"}

	for _, player := range r.Players {
		player.Send(message)
	}
}
func (r *Room) HandleSubmitLine(playerName, line string) {
	// r.Mutex.Lock()
	// defer r.Mutex.Unlock()
	if r.Status != "in_progress" {
		return
	}
	if r.TurnOrder[r.CurrentTurn] != playerName {
		return
	}
	r.Story = append(r.Story, line)
	r.BroadcastMessage(playerName + " added a line to the story. \nNew story: " + r.GetStory())
	r.NextTurn()
}

func (r *Room) NextTurn() {
	r.CurrentTurn = (r.CurrentTurn + 1) % len(r.TurnOrder)

	if r.isGameOver() { // Example end-game logic
		r.EndGame()
		return
	}
	r.BroadcastTurn()
}

func (r *Room) isGameOver() bool {
	return r.CurrentTurn == 0 && len(r.Story) >= r.TotalPlayers*5
}

func (r *Room) EndGame() {
	r.Status = "completed"
	endMessage := Message{Type: "END_GAME", Content: "Game over!"}
	for _, player := range r.Players {
		player.Send(endMessage)
	}
}
