package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyAccessService struct {
	parties *repository.PartyQueryRepository
}

func NewPartyAccessService(
	parties *repository.PartyQueryRepository,
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

	party, err := s.parties.FindByID(
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

	party, err := s.parties.FindByID(
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
