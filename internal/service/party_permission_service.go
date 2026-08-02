package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func (s *PermissionService) RequireManageParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
		enum.PartyRoleAdmin,
	)
}

func (s *PermissionService) CanManageParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.RequireManageParty(
		ctx,
		partyID,
		userID,
	) == nil
}

func (s *PermissionService) RequireDeleteParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
	)
}

func (s *PermissionService) CanDeleteParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.RequireDeleteParty(
		ctx,
		partyID,
		userID,
	) == nil
}

func (s *PermissionService) RequirePublishParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
	)
}

func (s *PermissionService) CanPublishParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.RequirePublishParty(
		ctx,
		partyID,
		userID,
	) == nil
}

func (s *PermissionService) RequireSchedulePublication(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
	)
}

func (s *PermissionService) CanSchedulePublication(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.RequireSchedulePublication(
		ctx,
		partyID,
		userID,
	) == nil
}

func (s *PermissionService) RequireManageRefunds(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
		enum.PartyRoleAdmin,
		enum.PartyRoleRefunder,
	)
}

func (s *PermissionService) CanManageRefunds(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.RequireManageRefunds(
		ctx,
		partyID,
		userID,
	) == nil
}

func (s *PermissionService) RequireScanTickets(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) error {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		enum.PartyRoleOrganizer,
		enum.PartyRoleAdmin,
		enum.PartyRoleScanner,
	)
}

func (s *PermissionService) CanScanTickets(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.RequireScanTickets(
		ctx,
		partyID,
		userID,
	) == nil
}
