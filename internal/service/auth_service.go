package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	auth_service "github.com/reinp/event-platform/backend/internal/service/auth"
)

type AuthService struct {
	users        *repository.UserRepository
	registration *auth_service.RegistrationService
	sessions     *auth_service.SessionService
	tokens       *auth_service.TokenService
	verification *auth_service.VerificationService
	passwords    *auth_service.PasswordService
}

func NewAuthService(
	users *repository.UserRepository,
	registration *auth_service.RegistrationService,
	sessions *auth_service.SessionService,
	tokens *auth_service.TokenService,
	verification *auth_service.VerificationService,
	passwords *auth_service.PasswordService,
) *AuthService {

	return &AuthService{
		users:        users,
		registration: registration,
		sessions:     sessions,
		tokens:       tokens,
		verification: verification,
		passwords:    passwords,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	req auth_service.RegisterRequest,
) (*models.User, error) {

	return s.registration.Register(
		ctx,
		req,
	)
}

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	token string,
) (string, error) {

	return s.verification.VerifyEmail(
		ctx,
		token,
	)
}

func (s *AuthService) Login(
	ctx context.Context,
	identifier string,
	password string,
	userAgent string,
	ipAddress string,
) (*auth_service.TokenResponse, error) {

	return s.sessions.Login(
		ctx,
		identifier,
		password,
		userAgent,
		ipAddress,
	)
}

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (*auth_service.TokenResponse, error) {

	return s.tokens.Refresh(
		ctx,
		refreshToken,
	)
}

func (s *AuthService) ResendVerificationEmail(
	ctx context.Context,
	email string,
) error {

	return s.verification.ResendVerificationEmail(
		ctx,
		email,
	)
}

func (s *AuthService) Secret() string {
	return s.tokens.Secret()
}

func (s *AuthService) ValidateUser(
	ctx context.Context,
	userID uuid.UUID,
) (*models.User, error) {

	return s.users.FindByID(
		ctx,
		userID,
	)
}

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {

	return s.sessions.Logout(
		ctx,
		refreshToken,
	)
}

func (s *AuthService) ForgotPassword(
	ctx context.Context,
	identifier string,
) error {

	return s.passwords.ForgotPassword(
		ctx,
		identifier,
	)
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	token string,
	newPassword string,
) error {

	return s.passwords.ResetPassword(
		ctx,
		token,
		newPassword,
	)
}

func (s *AuthService) Sessions(
	ctx context.Context,
	userID uuid.UUID,
	currentFamily uuid.UUID,
) ([]auth_service.SessionResponse, error) {

	return s.sessions.Sessions(
		ctx,
		userID,
		currentFamily,
	)
}

func (s *AuthService) RevokeSession(
	ctx context.Context,
	userID uuid.UUID,
	familyID uuid.UUID,
) error {

	return s.sessions.RevokeSession(
		ctx,
		userID,
		familyID,
	)
}

func (s *AuthService) LogoutAll(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return s.sessions.LogoutAll(
		ctx,
		userID,
	)
}
