package helpers

import (
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/mail"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
	"github.com/reinp/event-platform/backend/internal/service"
)

func NewAuthService(
	db *gorm.DB,
) *service.AuthService {

	executor := database.NewGormExecutor(db)

	// repositories

	userRepository :=
		repository.NewUserRepository(
			executor,
		)

	emailVerificationRepository :=
		repository.NewEmailVerificationRepository(
			executor,
		)

	refreshTokenRepository :=
		repository.NewRefreshTokenRepository(
			executor,
		)

	passwordResetRepository :=
		repository.NewPasswordResetTokenRepository(
			executor,
		)

	// jwt

	jwtService :=
		auth.NewJWT(
			"test-secret",
			NewFakeClock(
				time.Now(),
			),
		)

	// mail

	smtpPort, err := strconv.Atoi(
		os.Getenv("SMTP_PORT"),
	)

	if err != nil {
		panic(err)
	}

	mailer := mail.NewMailer(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
	)

	emailService :=
		service.NewEmailService(
			mailer,
			"http://localhost:5173",
		)

	// password validator

	passwordValidator :=
		security.NewPasswordValidator()

	return service.NewAuthService(
		userRepository,
		emailVerificationRepository,
		refreshTokenRepository,
		passwordResetRepository,
		jwtService,
		emailService,
		passwordValidator,
		24*time.Hour*30,
	)
}

func NewAuthServiceWithEmailService(
	db *gorm.DB,
	emailService service.EmailSender,
) *service.AuthService {

	executor := database.NewGormExecutor(db)

	userRepository :=
		repository.NewUserRepository(
			executor,
		)

	emailVerificationRepository :=
		repository.NewEmailVerificationRepository(
			executor,
		)

	refreshTokenRepository :=
		repository.NewRefreshTokenRepository(
			executor,
		)

	passwordResetRepository :=
		repository.NewPasswordResetTokenRepository(
			executor,
		)

	jwtService :=
		auth.NewJWT(
			"test-secret",
			NewFakeClock(time.Now()),
		)

	passwordValidator :=
		security.NewPasswordValidator()

	return service.NewAuthService(
		userRepository,
		emailVerificationRepository,
		refreshTokenRepository,
		passwordResetRepository,
		jwtService,
		emailService,
		passwordValidator,
		24*time.Hour*30,
	)
}

func NewAuthServiceWithCapturingEmail(
	db *gorm.DB,
	emailService *CapturingEmailService,
) *service.AuthService {

	return NewAuthServiceWithEmailService(
		db,
		emailService,
	)
}
