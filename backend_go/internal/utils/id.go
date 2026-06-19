package utils

import "github.com/google/uuid"

// GenerateID — buat UUID v4 baru.
func GenerateID() string {
	return uuid.New().String()
}
