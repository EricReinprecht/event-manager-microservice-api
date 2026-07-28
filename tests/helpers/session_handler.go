package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/handlers"
)

func NewSessionHandler(
	db *gorm.DB,
) *handlers.SessionHandler {

	authService := NewAuthService(
		db,
	)

	return handlers.NewSessionHandler(
		authService,
	)
}
