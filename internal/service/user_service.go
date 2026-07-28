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
	userRepository          *repository.UserRepository
	refreshTokensRepository *repository.RefreshTokenRepository
	passwordValidator       *security.PasswordValidator
}

func NewUserService(
	userRepository *repository.UserRepository,
	refreshTokensRepository *repository.RefreshTokenRepository,
	passwordValidator *security.PasswordValidator,
) *UserService {

	return &UserService{
		userRepository:          userRepository,
		refreshTokensRepository: refreshTokensRepository,
		passwordValidator:       passwordValidator,
	}
}

func (s *UserService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.User, error) {

	return s.userRepository.FindByID(
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

	user, err := s.userRepository.FindByID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName

	err = s.userRepository.Update(
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

	user, err := s.userRepository.FindByID(
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

	if err := s.userRepository.Update(
		ctx,
		user,
	); err != nil {

		return err
	}

	// logout everywhere

	if err := s.refreshTokensRepository.RevokeAllByUser(
		ctx,
		user.ID,
	); err != nil {

		return err
	}

	return nil
}
