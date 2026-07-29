package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
)

type UserService struct {
	users             *repository.UserRepository
	refreshTokens     *repository.RefreshTokenRepository
	passwordValidator *security.PasswordValidator
	passwordResets    *repository.PasswordResetTokenRepository
	verifications     *repository.EmailVerificationRepository
}

func NewUserService(
	users *repository.UserRepository,
	refreshTokens *repository.RefreshTokenRepository,
	passwordValidator *security.PasswordValidator,
	passwordResets *repository.PasswordResetTokenRepository,
	verifications *repository.EmailVerificationRepository,
) *UserService {

	return &UserService{
		users:             users,
		refreshTokens:     refreshTokens,
		passwordValidator: passwordValidator,
		passwordResets:    passwordResets,
		verifications:     verifications,
	}
}

func (s *UserService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {

	return s.users.FindByID(
		ctx,
		id,
	)
}

func (s *UserService) CompleteProfile(
	ctx context.Context,
	userID uuid.UUID,
	firstName string,
	lastName string,
) (*models.User, error) {

	user, err := s.users.FindByID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.ProfileCompleted = true

	err = s.users.Update(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	oldPassword string,
	newPassword string,
) error {

	user, err := s.users.FindByID(
		ctx,
		userID,
	)

	if err != nil {
		return errors.New(
			"user not found",
		)
	}

	// verify old password

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(oldPassword),
	)

	if err != nil {
		return errors.New(
			"invalid old password",
		)
	}

	// validate new password

	if err := s.passwordValidator.Validate(
		newPassword,
		user.Username,
		user.Email,
	); err != nil {

		return err
	}

	// hash new password

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)

	// update password

	if err := s.users.Update(
		ctx,
		user,
	); err != nil {

		return err
	}

	// logout everywhere

	if err := s.refreshTokens.RevokeAllByUser(
		ctx,
		user.ID,
	); err != nil {

		return err
	}

	return nil
}

func (s *UserService) Delete(
	ctx context.Context,
	userID uuid.UUID,
) error {

	// find user first
	user, err := s.users.FindByID(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	// revoke all active sessions
	err = s.refreshTokens.RevokeAllByUser(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	// invalidate password reset tokens
	err = s.passwordResets.InvalidateForUser(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	// optionally remove email verification tokens
	err = s.verifications.DeleteByUser(
		userID,
	)

	if err != nil {
		return err
	}

	// soft delete user
	err = s.users.Delete(
		ctx,
		user,
	)

	if err != nil {
		return err
	}

	return nil
}
