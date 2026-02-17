package utils

import (
	"github.com/google/uuid"
)

// GenerateID genera un UUID único
func GenerateID() string {
	return uuid.New().String()
}
