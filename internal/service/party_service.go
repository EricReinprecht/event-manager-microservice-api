package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyService struct {
	parties          *repository.PartyRepository
	memberRepository *repository.PartyMemberRepository
}

func NewPartyService(
	repository *repository.PartyRepository,
	memberRepository *repository.PartyMemberRepository,
) *PartyService {

	return &PartyService{
		parties:          repository,
		memberRepository: memberRepository,
	}
}

func (s *PartyService) Create(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	err := s.parties.Create(
		ctx,
		party,
		imageIDs,
	)

	if err != nil {
		return err
	}

	member := &models.PartyMember{
		UserID:  party.OrganizerID,
		PartyID: party.ID,
		Role:    enum.RoleOrganizer,
	}

	err = s.memberRepository.Create(
		ctx,
		member,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PartyService) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	return s.parties.FindAll(ctx)
}

func (s *PartyService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.parties.FindByID(ctx, id)
}

func (s *PartyService) Update(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Update(
		ctx,
		party,
	)
}

func (s *PartyService) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return s.parties.Delete(
		ctx,
		party,
	)
}

func (s *PartyService) UpdateImages(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	return s.parties.UpdateImages(
		ctx,
		party,
		imageIDs,
	)
}

func (s *PartyService) FindOwnedParty(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.Party, error) {

	party, err := s.FindByID(
		ctx,
		partyID,
	)

	if err != nil {
		return nil, err
	}

	if party.OrganizerID != userID {
		return nil, errors.New("not allowed")
	}

	return party, nil
}
