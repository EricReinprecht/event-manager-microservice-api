package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyService struct {
	parties *repository.PartyRepository
}

func NewPartyService(
	parties *repository.PartyRepository,
) *PartyService {

	return &PartyService{
		parties: parties,
	}
}

func (s *PartyService) Create(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Create(
		ctx,
		party,
	)
}

func (s *PartyService) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	return s.parties.FindAll(ctx)
}

func (s *PartyService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.parties.FindByID(ctx, id)
}
