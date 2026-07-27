package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
)

type AuthService struct {
	users             *repository.UserRepository
	jwt               *auth.JWT
	verifications     *repository.EmailVerificationRepository
	emailService      *EmailService
	passwordValidator *security.PasswordValidator
}

func NewAuthService(
	users *repository.UserRepository,
	jwt *auth.JWT,
	verifications *repository.EmailVerificationRepository,
	emailService *EmailService,
	passwordValidator *security.PasswordValidator,
) *AuthService {

	return &AuthService{
		users:             users,
		jwt:               jwt,
		verifications:     verifications,
		emailService:      emailService,
		passwordValidator: passwordValidator,
	}
}

type RegisterRequest struct {
	Email    string
	Password string
	Username string
}

func (s *AuthService) Register(
	ctx context.Context,
	req RegisterRequest,
) (*models.User, error) {

	_, err := s.users.FindByEmail(
		ctx,
		req.Email,
	)

	if err == nil {
		return nil, errors.New("email already exists")
	}

	_, err = s.users.FindByUsername(
		ctx,
		req.Username,
	)

	if err == nil {
		return nil, errors.New("username already exists")
	}

	err = s.passwordValidator.Validate(
		req.Password,
		req.Username,
		req.Email,
	)

	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Username:     req.Username,
	}

	err = s.users.Create(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	token := auth.GenerateToken()

	verification := &models.EmailVerification{
		UserID: user.ID,

		Token: token,

		ExpiresAt: time.Now().Add(
			24 * time.Hour,
		),
	}

	err = s.verifications.Create(
		verification,
	)

	if err != nil {
		return nil, err
	}

	err = s.emailService.SendVerificationEmail(
		user.Email,
		user.Username,
		token,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.users.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(
		user.ID.String(),
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Secret() string {
	return s.jwt.Secret
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

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	token string,
) (string, error) {

	verification, err := s.verifications.FindByToken(
		token,
	)

	if err != nil {
		return "", errors.New("invalid verification token")
	}

	if verification.ExpiresAt.Before(time.Now()) {
		return "", errors.New("verification expired")
	}

	user, err := s.users.FindByID(
		ctx,
		verification.UserID,
	)

	if err != nil {
		return "", err
	}

	// already verified
	if user.VerifiedAt != nil {

		return s.jwt.Generate(
			user.ID.String(),
		)
	}

	now := time.Now()

	user.VerifiedAt = &now

	err = s.users.Update(
		ctx,
		user,
	)

	if err != nil {
		return "", err
	}

	// remove used verification token
	err = s.verifications.Delete(
		verification.ID,
	)

	if err != nil {
		return "", err
	}

	return s.jwt.Generate(
		user.ID.String(),
	)
}
