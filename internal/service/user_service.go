package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(
	userRepository *repository.UserRepository,
) *UserService {

	return &UserService{
		userRepository: userRepository,
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
