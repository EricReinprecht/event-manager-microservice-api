package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"

	partyRepository "github.com/reinp/event-platform/backend/internal/repository/party"
)

type PartyAccessService struct {
	parties *partyRepository.Facade
}

func NewPartyAccessService(
	parties *partyRepository.Facade,
) *PartyAccessService {

	return &PartyAccessService{
		parties: parties,
	}
}

func (s *PartyAccessService) RequireOwnership(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.Party, error) {

	party, err := s.parties.Repository.FindByID(
		ctx,
		partyID,
	)

	if err != nil {

		return nil, err
	}

	if party.OrganizerID != userID {

		return nil, appErrors.ErrForbidden
	}

	return party, nil
}

func (s *PartyAccessService) FindOwnedParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.Party, error) {

	party, err := s.parties.Repository.FindByID(
		ctx,
		partyID,
	)

	if err != nil {
		return nil, err
	}

	if party.OrganizerID != userID {

		return nil, errors.New(
			"not allowed",
		)
	}

	return party, nil
}
