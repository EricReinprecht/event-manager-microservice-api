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
	imageIDs []uuid.UUID,
) error {

	return s.parties.Create(
		ctx,
		party,
		imageIDs,
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

func (s *PartyService) Update(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Update(
		ctx,
		party,
	)
}

func (s *PartyService) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Delete(
		ctx,
		party,
	)
}

func (s *PartyService) UpdateImages(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	return s.parties.UpdateImages(
		ctx,
		party,
		imageIDs,
	)
}
