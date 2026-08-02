package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func (s *PermissionService) CanManageParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.HasPartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
		enum.PartyRoleAdmin,
	)
}

func (s *PermissionService) CanScanTickets(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.HasPartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
		enum.PartyRoleAdmin,
	)
}

func (s *PermissionService) CanRefund(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.HasPartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
		enum.PartyRoleAdmin,
		enum.PartyRoleRefunder,
	)
}
