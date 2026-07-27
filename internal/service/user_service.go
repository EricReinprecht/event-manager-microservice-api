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
