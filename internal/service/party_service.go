package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/helpers"
	"github.com/reinp/event-platform/backend/internal/mapper"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/validators"
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

	categoryIDs, err := helpers.ParseUUIDs(
		req.CategoryIDs,
	)

	if err != nil {
		return nil, err
	}

	validationErrors := helpers.ValidateParty(
		party.StartAt,
		party.EndAt,
		party.Latitude,
		party.Longitude,
		party.Timezone,
	)

	helpers.MergeValidationErrors(
		validationErrors,
		validators.ValidateCreateTicketCategories(
			req.TicketCategories,
		),
	)

	if len(validationErrors) > 0 {
		return nil,
			appErrors.NewValidationError(
				validationErrors,
			)
	}

	if err := s.crud.CreateRelations(
		ctx,
		party,
		categoryIDs,
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
	req dto.UpdatePartyRequest,
) (*dto.PartyResponse, error) {

	party, err := s.query.FindByID(
		ctx,
		id,
	)

	if err != nil {
		return nil, err
	}

	validationErrors := helpers.ValidateParty(
		req.StartAt,
		req.EndAt,
		req.Location.Latitude,
		req.Location.Longitude,
		req.Location.Timezone,
	)

	helpers.MergeValidationErrors(
		validationErrors,
		validators.ValidateUpdateTicketCategories(
			req.TicketCategories,
		),
	)

	if len(validationErrors) > 0 {
		return nil,
			appErrors.NewValidationError(
				validationErrors,
			)
	}

	mapper.ApplyPartyUpdate(
		party,
		req,
	)

	categoryIDs, err := helpers.ParseUUIDs(
		req.CategoryIDs,
	)

	if err != nil {
		return nil, err
	}

	ticketCategories := mapper.TicketCategories(
		req.TicketCategories,
		party.ID,
	)

	if err := s.crud.UpdateRelations(
		ctx,
		party,
		categoryIDs,
		req.ImageIDs,
		ticketCategories,
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
) error {

	party, err := s.query.FindByID(
		ctx,
		id,
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
	locationName string,
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
		locationName,
		sorts,
		page,
		limit,
	)

	if err != nil {

		return nil, 0, err
	}

	return mapper.PartyResponses(parties), total, nil
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
