// internal/storage/memory_storage.go
package storage

import (
	"errors"
	"storytelling-backend/internal/models"
)

var memoryStore = make(map[string]*models.Room)

type MemoryStorage struct{}

func (ms *MemoryStorage) SaveRoom(room *models.Room) error {
	memoryStore[room.ID] = room
	return nil
}

func (ms *MemoryStorage) GetRoom(roomID string) (*models.Room, error) {
	room, exists := memoryStore[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (ms *MemoryStorage) DeleteRoom(roomID string) error {
	delete(memoryStore, roomID)
	return nil
}
