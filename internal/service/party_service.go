package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyService struct {
	crud   *PartyCRUDService
	query  *PartyQueryService
	access *PartyAccessService
}

func NewPartyService(
	crud *PartyCRUDService,
	query *PartyQueryService,
	access *PartyAccessService,
) *PartyService {

	return &PartyService{
		crud:   crud,
		query:  query,
		access: access,
	}
}

func (s *PartyService) Create(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	return s.crud.Create(
		ctx,
		party,
		imageIDs,
	)
}

func (s *PartyService) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	return s.query.FindAll(
		ctx,
	)
}

func (s *PartyService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.query.FindByID(
		ctx,
		id,
	)
}

func (s *PartyService) Update(
	ctx context.Context,
	party *models.Party,
) error {

	return s.crud.Update(
		ctx,
		party,
	)
}

func (s *PartyService) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return s.crud.Delete(
		ctx,
		party,
	)
}

func (s *PartyService) UpdateImages(
	ctx context.Context,
	partyID uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	return s.crud.UpdateImages(
		ctx,
		partyID,
		imageIDs,
	)
}

func (s *PartyService) FindOwnedParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.Party, error) {

	return s.access.FindOwnedParty(
		ctx,
		partyID,
		userID,
	)
}

func (s *PartyService) FindForUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	role string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	return s.query.FindForUser(
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

func (s *PartyService) FindOrganizedByUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	sorts string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	return s.query.FindOrganizedByUser(
		ctx,
		userID,
		name,
		startAt,
		endAt,
		sorts,
		page,
		limit,
	)
}
