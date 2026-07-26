package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/mail"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type AuthService struct {
	users         *repository.UserRepository
	jwt           *auth.JWT
	mailer        *mail.Mailer
	verifications *repository.EmailVerificationRepository
}

func NewAuthService(
	users *repository.UserRepository,
	jwt *auth.JWT,
	mailer *mail.Mailer,
	verifications *repository.EmailVerificationRepository,
) *AuthService {

	return &AuthService{
		users:         users,
		jwt:           jwt,
		mailer:        mailer,
		verifications: verifications,
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

	err = s.users.Create(ctx, user)

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
