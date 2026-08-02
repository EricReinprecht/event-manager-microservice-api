package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
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

	err := s.repository.Create(
		s.repository.DB(ctx),
		member,
	)

	if err != nil {
		return nil, mapPartyMemberDatabaseError(err)
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

	return s.repository.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			roleRepository :=
				repository.NewPartyMemberRoleRepository(tx)

			if err := roleRepository.DeleteAll(
				tx,
				memberID,
			); err != nil {
				return err
			}

			for _, role := range roles {

				memberRole := &models.PartyMemberRole{
					ID:            uuid.New(),
					PartyMemberID: memberID,
					Role:          role,
				}

				if err := roleRepository.Create(
					tx,
					memberRole,
				); err != nil {
					return mapPartyMemberDatabaseError(err)
				}
			}

			return nil
		},
	)
}

func mapPartyMemberDatabaseError(
	err error,
) error {

	var pgErr *pgconn.PgError

	if !errors.As(
		err,
		&pgErr,
	) {
		return err
	}

	switch {

	case pgErr.Code == "23505" &&
		pgErr.ConstraintName == "idx_party_member_user_party":

		return appErrors.ErrPartyMemberAlreadyExists

	case pgErr.Code == "23514":

		return appErrors.ErrInvalidPartyMemberRole

	default:

		return err
	}
}
