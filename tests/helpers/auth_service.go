package helpers

import (
	"time"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewAuthService(
	db *gorm.DB,
) *service.AuthService {

	userRepository := repository.NewUserRepository(
		db,
	)

	jwt := auth.NewJWT(
		"test-secret",
		NewFakeClock(
			time.Now(),
		),
	)

	return service.NewAuthService(
		userRepository,
		jwt,
	)
}
