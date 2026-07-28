package scenarios

import (
	"context"

	"github.com/google/uuid"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/fixtures"
)

type AuthScenario struct {
	DB         *gorm.DB
	Service    *service.AuthService
	UserRepo   *repository.UserRepository
	TokenRepo  *repository.RefreshTokenRepository
	ResetRepo  *repository.PasswordResetTokenRepository
	VerifyRepo *repository.EmailVerificationRepository
}

func NewAuthScenario(
	db *gorm.DB,
	authService *service.AuthService,
) *AuthScenario {

	executor := database.NewGormExecutor(
		db,
	)

	return &AuthScenario{

		DB: db,

		Service: authService,

		UserRepo: repository.NewUserRepository(
			executor,
		),

		TokenRepo: repository.NewRefreshTokenRepository(
			executor,
		),

		ResetRepo: repository.NewPasswordResetTokenRepository(
			executor,
		),

		VerifyRepo: repository.NewEmailVerificationRepository(
			executor,
		),
	}
}

// Create verified user ready for login tests
func (s *AuthScenario) CreateVerifiedUser(
	ctx context.Context,
) *models.User {

	user := fixtures.VerifiedUser()

	err := s.DB.WithContext(ctx).
		Create(&user).
		Error

	if err != nil {
		panic(err)
	}

	return &user
}

// Create normal unverified user
func (s *AuthScenario) CreateUser(
	ctx context.Context,
) *models.User {

	user := fixtures.User()

	err := s.DB.WithContext(ctx).
		Create(&user).
		Error

	if err != nil {
		panic(err)
	}

	return &user
}

// Login user and create refresh session
func (s *AuthScenario) Login(
	ctx context.Context,
	user *models.User,
) *service.TokenResponse {

	response, err := s.Service.Login(
		ctx,
		user.Email,
		"C4ctus!River#829Lamp",
		"test-agent",
		"127.0.0.1",
	)

	if err != nil {
		panic(err)
	}

	return response
}

// Create a fake refresh session
func (s *AuthScenario) CreateRefreshToken(
	ctx context.Context,
	userID uuid.UUID,
) *models.RefreshToken {

	token := fixtures.RefreshToken(
		userID,
	)

	err := s.DB.WithContext(ctx).
		Create(&token).
		Error

	if err != nil {
		panic(err)
	}

	return &token
}

// Create email verification token
func (s *AuthScenario) CreateVerification(
	ctx context.Context,
	userID uuid.UUID,
	rawToken string,
) *models.EmailVerification {

	verification := fixtures.EmailVerification(
		userID,
		auth.HashToken(rawToken),
	)

	err := s.DB.WithContext(ctx).
		Create(&verification).
		Error

	if err != nil {
		panic(err)
	}

	return &verification
}

// Create password reset token
func (s *AuthScenario) CreatePasswordReset(
	ctx context.Context,
	userID uuid.UUID,
	rawToken string,
) *models.PasswordResetToken {

	reset := fixtures.PasswordResetToken(
		userID,
		auth.HashToken(rawToken),
	)

	err := s.DB.WithContext(ctx).
		Create(&reset).
		Error

	if err != nil {
		panic(err)
	}

	return &reset
}
