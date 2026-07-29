package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

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
