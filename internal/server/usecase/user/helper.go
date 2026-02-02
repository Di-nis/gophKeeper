package user

import (
	"github.com/google/uuid"
)

// generateSessionID - генерация уникального идентификатора сессии.
func generateSessionID() string {
	return uuid.NewString()
}
