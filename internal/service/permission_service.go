package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PermissionService struct {
	partyService       *PartyService
	partyMemberService *PartyMemberService
}

func NewPermissionService(
	partyService *PartyService,
	partyMemberService *PartyMemberService,
) *PermissionService {

	return &PermissionService{
		partyService:       partyService,
		partyMemberService: partyMemberService,
	}
}

func (s *PermissionService) RequirePartyRole(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
	roles ...enum.PartyRole,
) error {

	party, err := s.partyService.FindByID(
		ctx,
		partyID,
	)

	if err == nil {

		if party.OrganizerID == userID {

			for _, role := range roles {

				if role == enum.RoleOrganizer {
					return nil
				}
			}
		}
	}

	member, err := s.partyMemberService.FindByPartyAndUser(
		ctx,
		partyID,
		userID,
	)

	if err != nil {
		return appErrors.ErrNotAllowed
	}

	for _, wanted := range roles {

		if member.HasRole(wanted) {
			return nil
		}
	}

	return appErrors.ErrNotAllowed
}

func (s *PermissionService) HasPartyRole(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
	roles ...enum.PartyRole,
) bool {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		roles...,
	) == nil
}
