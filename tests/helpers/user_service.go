package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewUserService(
	db *gorm.DB,
) *service.UserService {

	executor := database.NewGormExecutor(
		db,
	)

	userRepository := repository.NewUserRepository(
		executor,
	)

	refreshTokenRepository := repository.NewRefreshTokenRepository(
		executor,
	)

	passwordValidator := security.NewPasswordValidator()

	return service.NewUserService(
		userRepository,
		refreshTokenRepository,
		passwordValidator,
	)
}
