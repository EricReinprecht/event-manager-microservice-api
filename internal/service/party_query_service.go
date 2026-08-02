package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"

	partyRepository "github.com/reinp/event-platform/backend/internal/repository/party"
)

type PartyQueryService struct {
	parties *partyRepository.Facade
}

func NewPartyQueryService(
	parties *partyRepository.Facade,
) *PartyQueryService {

	return &PartyQueryService{
		parties: parties,
	}
}

func (s *PartyQueryService) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	return s.parties.Repository.FindAll(
		ctx,
	)
}

func (s *PartyQueryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.parties.Repository.FindByID(
		ctx,
		id,
	)
}

func (s *PartyQueryService) FindForUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	role string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	return s.parties.Repository.FindForUser(
		ctx,
		userID,
		name,
		startAt,
		endAt,
		role,
		page,
		limit,
	)
}

func (s *PartyQueryService) FindOrganizedByUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	locationName string,
	sorts string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	return s.parties.Repository.FindOrganizedByUser(
		ctx,
		userID,
		name,
		startAt,
		endAt,
		locationName,
		sorts,
		page,
		limit,
	)
}
