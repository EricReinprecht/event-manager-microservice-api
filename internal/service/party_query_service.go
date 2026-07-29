package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyQueryService struct {
	parties *repository.PartyQueryRepository
}

func NewPartyQueryService(
	parties *repository.PartyQueryRepository,
) *PartyQueryService {

	return &PartyQueryService{
		parties: parties,
	}
}

func (s *PartyQueryService) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	return s.parties.FindAll(ctx)
}

func (s *PartyQueryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	return s.parties.FindByID(
		ctx,
		id,
	)
}

func (s *PartyQueryService) FindForUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	role string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	return s.parties.FindForUser(
		ctx,
		userID,
		name,
		startAt,
		endAt,
		role,
		page,
		limit,
	)
}

func (s *PartyQueryService) FindOrganizedByUser(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	startAt string,
	endAt string,
	sorts string,
	page int,
	limit int,
) ([]models.Party, int64, error) {

	return s.parties.FindOrganizedByUser(
		ctx,
		userID,
		name,
		startAt,
		endAt,
		sorts,
		page,
		limit,
	)
}
