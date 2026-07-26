package helpers

import (
	"time"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/mail"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewAuthService(
	db *gorm.DB,
) *service.AuthService {

	executor := database.NewGormExecutor(db)

	userRepository := repository.NewUserRepository(
		executor,
	)

	emailVerificationRepository := repository.NewEmailVerificationRepository(
		executor,
	)

	jwt := auth.NewJWT(
		"test-secret",
		NewFakeClock(
			time.Now(),
		),
	)

	mailer := mail.NewMailer(
		"localhost",
		1025,
		"test",
		"test",
		"test@example.com",
	)

	return service.NewAuthService(
		userRepository,
		jwt,
		mailer,
		emailVerificationRepository,
	)
}
