package auth_service

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
)

type PasswordService struct {
	users          *repository.UserRepository
	passwordResets *repository.PasswordResetTokenRepository
	refreshTokens  *repository.RefreshTokenRepository

	passwordValidator *security.PasswordValidator
	emailService      EmailSender
	clock             clock.Clock

	resetDuration time.Duration
	resetCooldown time.Duration
}

func NewPasswordService(
	users *repository.UserRepository,
	passwordResets *repository.PasswordResetTokenRepository,
	refreshTokens *repository.RefreshTokenRepository,
	passwordValidator *security.PasswordValidator,
	emailService EmailSender,
	clock clock.Clock,
	resetDuration time.Duration,
	resetCooldown time.Duration,
) *PasswordService {

	return &PasswordService{
		users:             users,
		passwordResets:    passwordResets,
		refreshTokens:     refreshTokens,
		passwordValidator: passwordValidator,
		emailService:      emailService,
		clock:             clock,
		resetDuration:     resetDuration,
		resetCooldown:     resetCooldown,
	}
}

func (s *PasswordService) now() time.Time {
	return s.clock.Now()
}

func (s *PasswordService) ForgotPassword(
	ctx context.Context,
	identifier string,
) error {

	identifier = strings.TrimSpace(identifier)

	user, err := s.users.FindByIdentifier(
		ctx,
		identifier,
	)

	if err != nil {
		// Do not reveal whether the account exists.
		return nil
	}

	latestToken, err := s.passwordResets.FindLatestByUser(
		ctx,
		user.ID,
	)

	if err == nil &&
		s.now().Sub(latestToken.CreatedAt) < s.resetCooldown {

		// Silently ignore requests during cooldown.
		return nil
	}

	if err := s.passwordResets.InvalidateForUser(
		ctx,
		user.ID,
	); err != nil {
		return err
	}

	rawToken := auth.GenerateToken()

	resetToken := &models.PasswordResetToken{
		UserID: user.ID,

		TokenHash: auth.HashToken(
			rawToken,
		),

		ExpiresAt: s.now().Add(
			s.resetDuration,
		),
	}

	if err := s.passwordResets.Create(
		ctx,
		resetToken,
	); err != nil {
		return err
	}

	return s.emailService.SendPasswordResetEmail(
		user.Email,
		user.Username,
		rawToken,
	)
}

func (s *PasswordService) ResetPassword(
	ctx context.Context,
	token string,
	newPassword string,
) error {

	tokenHash := auth.HashToken(token)

	resetToken, err := s.passwordResets.FindByHash(
		ctx,
		tokenHash,
	)

	if err != nil {
		return errors.New(
			"invalid reset token",
		)
	}

	if resetToken.InvalidatedAt != nil {
		return errors.New(
			"reset token invalidated",
		)
	}

	if resetToken.UsedAt != nil {
		return errors.New(
			"reset token already used",
		)
	}

	if resetToken.ExpiresAt.Before(s.now()) {
		return errors.New(
			"reset token expired",
		)
	}

	user, err := s.users.FindByID(
		ctx,
		resetToken.UserID,
	)

	if err != nil {
		return errors.New(
			"cannot reset password for deleted user",
		)
	}

	if err := s.passwordValidator.Validate(
		newPassword,
		user.Username,
		user.Email,
	); err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)

	if err := s.users.Update(
		ctx,
		user,
	); err != nil {
		return err
	}

	if err := s.passwordResets.MarkUsed(
		ctx,
		resetToken.ID,
	); err != nil {
		return err
	}

	return s.refreshTokens.RevokeAllByUser(
		ctx,
		user.ID,
	)
}
