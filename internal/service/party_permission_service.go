package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

func (s *PartyMemberService) CanManageParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.HasRole(
		ctx,
		partyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
	)
}

func (s *PartyMemberService) CanScanTickets(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.HasRole(
		ctx,
		partyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
		enum.RoleStaff,
	)
}

func (s *PartyMemberService) CanRefund(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) bool {

	return s.HasRole(
		ctx,
		partyID,
		userID,
		enum.RoleOrganizer,
		enum.RoleAdmin,
		enum.RoleRefunder,
	)
}
