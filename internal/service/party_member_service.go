package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyMemberService struct {
	repository      *repository.PartyMemberRepository
	partyRepository *repository.PartyRepository
}

func NewPartyMemberService(
	repository *repository.PartyMemberRepository,
	partyRepository *repository.PartyRepository,
) *PartyMemberService {

	return &PartyMemberService{
		repository:      repository,
		partyRepository: partyRepository,
	}
}

func (s *PartyMemberService) Create(
	ctx context.Context,
	partyID uuid.UUID,
	req dto.CreatePartyMemberRequest,
) (*models.PartyMember, error) {

	member := &models.PartyMember{
		ID:      uuid.New(),
		UserID:  req.UserID,
		PartyID: partyID,
		Roles: make(
			[]models.PartyMemberRole,
			0,
			len(req.Roles),
		),
	}

	for _, role := range req.Roles {

		member.Roles = append(
			member.Roles,
			models.PartyMemberRole{
				ID:            uuid.New(),
				PartyMemberID: member.ID,
				Role:          role,
			},
		)
	}

	if err := s.repository.Create(
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

	return s.repository.FindByID(
		ctx,
		id,
	)
}

func (s *PartyMemberService) FindByPartyAndUser(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.PartyMember, error) {

	return s.repository.FindByPartyAndUser(
		ctx,
		partyID,
		userID,
	)
}

func (s *PartyMemberService) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.PartyMember, error) {

	return s.repository.FindByParty(
		ctx,
		partyID,
	)
}

func (s *PartyMemberService) Update(
	ctx context.Context,
	member *models.PartyMember,
) error {

	return s.repository.Update(
		ctx,
		member,
	)
}

func (s *PartyMemberService) Delete(
	ctx context.Context,
	memberID uuid.UUID,
) error {

	member, err := s.repository.FindByID(
		ctx,
		memberID,
	)

	if err != nil {
		return err
	}

	party, err := s.partyRepository.FindByID(
		ctx,
		member.PartyID,
	)

	if err != nil {
		return err
	}

	if party.OrganizerID == member.UserID {
		return appErrors.ErrCannotRemoveOrganizer
	}

	return s.repository.Delete(
		ctx,
		member,
	)
}

func (s *PartyMemberService) SyncRoles(
	ctx context.Context,
	memberID uuid.UUID,
	roles []enum.PartyMemberRole,
) error {

	if _, err := s.repository.FindByID(
		ctx,
		memberID,
	); err != nil {
		return err
	}

	if len(roles) == 0 {
		return appErrors.ErrInvalidPartyMemberRole
	}

	uniqueRoles := make(
		[]enum.PartyMemberRole,
		0,
		len(roles),
	)

	seenRoles := make(
		map[enum.PartyMemberRole]struct{},
		len(roles),
	)

	for _, role := range roles {

		if _, exists := seenRoles[role]; exists {
			continue
		}

		seenRoles[role] = struct{}{}

		uniqueRoles = append(
			uniqueRoles,
			role,
		)
	}

	return s.repository.SyncRoles(
		ctx,
		memberID,
		uniqueRoles,
	)
}
