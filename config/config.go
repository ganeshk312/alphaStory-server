// config/config.go
package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Loading default configuration.")
	}
}
