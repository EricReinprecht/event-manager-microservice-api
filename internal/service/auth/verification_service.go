package auth_service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type VerificationService struct {
	users         *repository.UserRepository
	verifications *repository.EmailVerificationRepository
	tokens        *TokenService
	clock         clock.Clock
	emailService  EmailSender

	duration time.Duration
	cooldown time.Duration
}

func NewVerificationService(
	users *repository.UserRepository,
	verifications *repository.EmailVerificationRepository,
	tokens *TokenService,
	clock clock.Clock,
	emailService EmailSender,
	duration time.Duration,
	cooldown time.Duration,
) *VerificationService {

	return &VerificationService{
		users:         users,
		verifications: verifications,
		tokens:        tokens,
		clock:         clock,
		emailService:  emailService,
		duration:      duration,
		cooldown:      cooldown,
	}
}

func (s *VerificationService) now() time.Time {
	return s.clock.Now()
}

func (s *VerificationService) NewVerification(
	userID uuid.UUID,
) (*models.EmailVerification, string) {

	rawToken := auth.GenerateToken()

	verification := &models.EmailVerification{
		UserID: userID,

		Token: auth.HashToken(
			rawToken,
		),

		ExpiresAt: s.now().Add(
			s.duration,
		),
	}

	return verification, rawToken
}

func (s *VerificationService) VerifyEmail(
	ctx context.Context,
	token string,
) (string, error) {

	verification, err := s.verifications.FindByToken(
		ctx,
		auth.HashToken(token),
	)

	if err != nil {
		return "", errors.New(
			"invalid verification token",
		)
	}

	if verification.ExpiresAt.Before(s.now()) {
		return "", errors.New(
			"verification expired",
		)
	}

	user, err := s.users.FindByID(
		ctx,
		verification.UserID,
	)

	if err != nil {
		return "", err
	}

	if user.VerifiedAt == nil {

		now := s.now()

		user.VerifiedAt = &now

		if err := s.users.Update(
			ctx,
			user,
		); err != nil {
			return "", err
		}

		if err := s.verifications.Delete(
			verification.ID,
		); err != nil {
			return "", err
		}
	}

	familyID := uuid.New()

	return s.tokens.GenerateAccessToken(
		user.ID,
		familyID,
	)
}

func (s *VerificationService) ResendVerificationEmail(
	ctx context.Context,
	email string,
) error {

	email = strings.TrimSpace(
		strings.ToLower(email),
	)

	user, err := s.users.FindByEmail(
		ctx,
		email,
	)

	if err != nil || user.DeletedAt.Valid {
		return errors.New(
			"invalid request",
		)
	}

	if user.VerifiedAt != nil {
		return errors.New(
			"email already verified",
		)
	}

	lastVerification, err :=
		s.verifications.FindLatestByUser(
			ctx,
			user.ID,
		)

	if err == nil &&
		lastVerification.CreatedAt.After(
			s.now().Add(-s.cooldown),
		) {

		return errors.New(
			"please wait before requesting another verification email",
		)
	}

	if err := s.verifications.InvalidateForUser(
		ctx,
		user.ID,
	); err != nil {
		return err
	}

	verification, rawToken :=
		s.NewVerification(user.ID)

	if err := s.verifications.Create(
		verification,
	); err != nil {
		return err
	}

	return s.emailService.SendVerificationEmail(
		user.Email,
		user.Username,
		rawToken,
	)
}
