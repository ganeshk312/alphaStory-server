// pkg/utils/utils.go
package utils

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/exp/rand"
)

func LogInfo(message string) {
	log.Println("INFO: " + message)
}

func LogError(err error) {
	if err != nil {
		log.Println("ERROR: " + err.Error())
	}
}

// GenerateRoomID creates a unique random room ID
func GenerateRoomID() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	return fmt.Sprintf("rm-%d-%d", rand.Intn(1000), time.Now().UnixMicro()%1000)
}
