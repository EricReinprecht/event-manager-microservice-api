package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	appErrors "github.com/reinp/event-platform/backend/internal/apperrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type PartyMemberService struct {
	repository *repository.PartyMemberRepository
}

func NewPartyMemberService(
	repository *repository.PartyMemberRepository,
) *PartyMemberService {

	return &PartyMemberService{
		repository: repository,
	}
}

func (s *PartyMemberService) Create(
	ctx context.Context,
	member *models.PartyMember,
) error {

	err := s.repository.Create(ctx, member)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "idx_party_member_user_party" {

				return appErrors.ErrPartyMemberAlreadyExists
			}
		}

		return err
	}

	return nil
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
	member *models.PartyMember,
) error {

	return s.repository.Delete(
		ctx,
		member,
	)
}

// Permission helper
func (s *PartyMemberService) HasRole(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
	roles ...enum.PartyRole,
) bool {

	member, err := s.repository.FindByPartyAndUser(
		ctx,
		partyID,
		userID,
	)

	if err != nil {
		return false
	}

	for _, role := range roles {

		if member.Role == role {
			return true
		}
	}

	return false
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
