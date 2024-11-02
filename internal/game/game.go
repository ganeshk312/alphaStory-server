// internal/game/room_manager.go
package game

import (
	"errors"
	"storytelling-backend/internal/models"
	"sync"
)

// RoomManager is responsible for managing rooms and players.
type RoomManager struct {
	rooms      map[string]*models.Room
	roomsMutex sync.RWMutex
}

var RoomManagerInstance *RoomManager

// NewRoomManager creates and returns a new RoomManager.
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*models.Room),
	}
}

// CreateRoom creates a new room and adds it to the manager.
func (rm *RoomManager) CreateRoom(roomID string, host string) (*models.Room, error) {
	// rm.roomsMutex.Lock()
	// defer rm.roomsMutex.Unlock()

	if _, exists := rm.rooms[roomID]; exists {
		return nil, errors.New("room already exists")
	}

	room := models.NewRoom(roomID, host)
	rm.rooms[roomID] = room
	return room, nil
}

// GetRoom retrieves a room by ID.
func (rm *RoomManager) GetRoom(roomID string) (*models.Room, error) {
	// rm.roomsMutex.RLock()
	// defer rm.roomsMutex.RUnlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	return room, nil
}

// AddPlayerToRoom adds a player to the specified room.
func (rm *RoomManager) AddPlayerToRoom(roomID, playerName string) (*models.Room, error) {
	// rm.roomsMutex.Lock()
	// defer rm.roomsMutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	err := room.AddPlayer(playerName)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// AddConnectionToRoom adds a WebSocket connection for a player in a specific room.
func (rm *RoomManager) AddConnectionToRoom(roomID string, conn *models.PlayerConnection) error {
	// rm.roomsMutex.Lock()
	// defer rm.roomsMutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	err := room.AddConnection(conn)
	if err != nil {
		return err
	}
	return nil
}

func (rm *RoomManager) GetStory(roomID string) ([]string, error) {
	// rm.roomsMutex.Lock()
	// defer rm.roomsMutex.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	return room.Story, nil
}

// // BroadcastStoryUpdate broadcasts the updated story to all connected players in a room.
// func (rm *RoomManager) BroadcastStoryUpdate(roomID string) error {
// rm.roomsMutex.RLock()
// defer rm.roomsMutex.RUnlock()

// 	room, exists := rm.rooms[roomID]
// 	if !exists {
// 		return errors.New("room not found")
// 	}

// 	room.BroadcastStoryUpdate()
// 	return nil
// }
