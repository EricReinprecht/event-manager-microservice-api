package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PermissionService struct {
	partyRepository       *repository.PartyRepository
	partyMemberRepository *repository.PartyMemberRepository
}

func NewPermissionService(
	partyRepository *repository.PartyRepository,
	partyMemberRepository *repository.PartyMemberRepository,
) *PermissionService {

	return &PermissionService{
		partyRepository:       partyRepository,
		partyMemberRepository: partyMemberRepository,
	}
}

func (s *PermissionService) RequirePartyRole(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
	roles ...enum.PartyMemberRole,
) error {

	if len(roles) == 0 {
		return appErrors.ErrNotAllowed
	}

	party, err := s.partyRepository.FindByID(
		ctx,
		partyID,
	)

	if err != nil {
		return err
	}

	if party.OrganizerID == userID &&
		containsPartyRole(
			roles,
			enum.PartyRoleOrganizer,
		) {

		return nil
	}

	member, err :=
		s.partyMemberRepository.FindByPartyAndUser(
			ctx,
			partyID,
			userID,
		)

	if err != nil {
		return appErrors.ErrNotAllowed
	}

	for _, role := range roles {

		if member.HasRole(role) {
			return nil
		}
	}

	return appErrors.ErrNotAllowed
}

func (s *PermissionService) HasPartyRole(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
	roles ...enum.PartyMemberRole,
) bool {

	return s.RequirePartyRole(
		ctx,
		partyID,
		userID,
		roles...,
	) == nil
}

func containsPartyRole(
	roles []enum.PartyMemberRole,
	wanted enum.PartyMemberRole,
) bool {

	for _, role := range roles {

		if role == wanted {
			return true
		}
	}

	return false
}
