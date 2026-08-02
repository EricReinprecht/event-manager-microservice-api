package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/mapper"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	partyRepository "github.com/reinp/event-platform/backend/internal/repository/party"
	partyMemberRepository "github.com/reinp/event-platform/backend/internal/repository/party_member"
)

type PartyMemberService struct {
	members *partyMemberRepository.Facade
	parties *partyRepository.Facade
}

func NewPartyMemberService(
	members *partyMemberRepository.Facade,
	parties *partyRepository.Facade,
) *PartyMemberService {

	return &PartyMemberService{
		members: members,
		parties: parties,
	}
}

func (s *PartyMemberService) Create(
	ctx context.Context,
	partyID uuid.UUID,
	req dto.CreatePartyMemberRequest,
) (*models.PartyMember, error) {

	member := mapper.CreatePartyMember(
		partyID,
		req,
	)

	if err := s.members.Repository.Create(
		ctx,
		member,
	); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *PartyMemberService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.PartyMember, error) {

	return s.members.Repository.FindByID(
		ctx,
		id,
	)
}

func (s *PartyMemberService) FindByPartyAndUser(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.PartyMember, error) {

	return s.members.Repository.FindByPartyAndUser(
		ctx,
		partyID,
		userID,
	)
}

func (s *PartyMemberService) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.PartyMember, error) {

	return s.members.Repository.FindByParty(
		ctx,
		partyID,
	)
}

func (s *PartyMemberService) Update(
	ctx context.Context,
	member *models.PartyMember,
) error {

	return s.members.Repository.Update(
		ctx,
		member,
	)
}

func (s *PartyMemberService) Delete(
	ctx context.Context,
	memberID uuid.UUID,
) error {

	member, err := s.members.Repository.FindByID(
		ctx,
		memberID,
	)

	if err != nil {
		return err
	}

	party, err := s.parties.Repository.FindByID(
		ctx,
		member.PartyID,
	)

	if err != nil {
		return err
	}

	if party.OrganizerID == member.UserID {
		return appErrors.ErrCannotRemoveOrganizer
	}

	return s.members.Repository.Delete(
		ctx,
		member,
	)
}

func (s *PartyMemberService) SyncRoles(
	ctx context.Context,
	memberID uuid.UUID,
	roles []enum.PartyMemberRole,
) error {

	if _, err := s.members.Repository.FindByID(
		ctx,
		memberID,
	); err != nil {

		return err
	}

	if len(roles) == 0 {
		return appErrors.ErrInvalidPartyMemberRole
	}

	return s.members.Write.SyncRoles(
		ctx,
		memberID,
		mapper.UniquePartyRoles(
			roles,
		),
	)
}
