package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/mapper"
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
	req dto.CreatePartyRequest,
	userID uuid.UUID,
) (*dto.PartyResponse, error) {

	party := mapper.NewParty(
		req,
		userID,
	)

	if err := helpers.ValidateParty(
		party.StartAt,
		party.EndAt,
		party.Latitude,
		party.Longitude,
		party.Timezone,
	); err != nil {

		return nil, err
	}

	if err := s.crud.Create(
		ctx,
		party,
		req.ImageIDs,
	); err != nil {

		return nil, err
	}

	createdParty, err := s.query.FindByID(
		ctx,
		party.ID,
	)

	if err != nil {

		return nil, err
	}

	response := mapper.PartyResponse(
		createdParty,
	)

	return &response, nil
}

func (s *PartyService) FindAll(
	ctx context.Context,
) ([]dto.PartyResponse, error) {

	parties, err := s.query.FindAll(
		ctx,
	)

	if err != nil {

		return nil, err
	}

	return mapper.PartyResponses(parties), nil
}

func (s *PartyService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.PartyResponse, error) {

	party, err := s.query.FindByID(
		ctx,
		id,
	)

	if err != nil {

		return nil, err
	}

	response := mapper.PartyResponse(
		party,
	)

	return &response, nil
}

func (s *PartyService) Update(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	req dto.UpdatePartyRequest,
) (*dto.PartyResponse, error) {

	party, err := s.access.RequireOwnership(
		ctx,
		id,
		userID,
	)

	if err != nil {

		return nil, err
	}

	if err := helpers.ValidateParty(
		req.StartAt,
		req.EndAt,
		req.Latitude,
		req.Longitude,
		req.Timezone,
	); err != nil {

		return nil, err
	}

	mapper.ApplyPartyUpdate(
		party,
		req,
	)

	if err := s.crud.Update(
		ctx,
		party,
	); err != nil {

		return nil, err
	}

	if err := s.crud.UpdateImages(
		ctx,
		party.ID,
		req.ImageIDs,
	); err != nil {

		return nil, err
	}

	updatedParty, err := s.query.FindByID(
		ctx,
		id,
	)

	if err != nil {

		return nil, err
	}

	response := mapper.PartyResponse(
		updatedParty,
	)

	return &response, nil
}

func (s *PartyService) Delete(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) error {

	party, err := s.access.RequireOwnership(
		ctx,
		id,
		userID,
	)

	if err != nil {

		return err
	}

	return s.crud.Delete(
		ctx,
		party,
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
) ([]dto.PartyResponse, int64, error) {

	parties, total, err := s.query.FindOrganizedByUser(
		ctx,
		userID,
		name,
		startAt,
		endAt,
		sorts,
		page,
		limit,
	)

	if err != nil {

		return nil, 0, err
	}

	return mapper.PartyResponses(parties), total, nil
}
