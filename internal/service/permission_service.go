package service

import (
	"context"

	"github.com/google/uuid"
)

type PermissionService struct {
	partyService *PartyService
}

func NewPermissionService(
	partyService *PartyService,
) *PermissionService {

	return &PermissionService{
		partyService: partyService,
	}
}

func (s *PermissionService) RequirePartyOwner(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	_, err := s.partyService.FindOwnedParty(
		ctx,
		partyID,
		userID,
	)

	return err
}
