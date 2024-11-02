// internal/storage/storage.go
package storage

import "storytelling-backend/internal/models"

type Storage interface {
	SaveRoom(roomID string) error
	GetRoom(roomID string) (*models.Room, error)
	DeleteRoom(roomID string) error
}
