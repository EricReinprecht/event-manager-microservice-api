package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMemberRepository struct {
	db                 database.DBExecutor
	transactionManager *database.TransactionManager
}

func NewPartyMemberRepository(
	db database.DBExecutor,
	transactionManager *database.TransactionManager,
) *PartyMemberRepository {

	return &PartyMemberRepository{
		db:                 db,
		transactionManager: transactionManager,
	}
}

func (r *PartyMemberRepository) Create(
	ctx context.Context,
	member *models.PartyMember,
) error {

	err := r.db.
		WithContext(ctx).
		Create(member).
		Error()

	if err != nil {
		return mapPartyMemberDatabaseError(err)
	}

	return nil
}

func (r *PartyMemberRepository) FindByPartyAndUser(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.PartyMember, error) {

	var member models.PartyMember

	err := r.db.
		WithContext(ctx).
		Preload("Roles").
		Where(
			"party_id = ? AND user_id = ?",
			partyID,
			userID,
		).
		First(&member).
		Error()

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *PartyMemberRepository) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.PartyMember, error) {

	var members []models.PartyMember

	err := r.db.
		WithContext(ctx).
		Preload("Roles").
		Where(
			"party_id = ?",
			partyID,
		).
		Find(&members).
		Error()

	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *PartyMemberRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.PartyMember, error) {

	var member models.PartyMember

	err := r.db.
		WithContext(ctx).
		Preload("Roles").
		First(
			&member,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *PartyMemberRepository) Update(
	ctx context.Context,
	member *models.PartyMember,
) error {

	err := r.db.
		WithContext(ctx).
		Save(member).
		Error()

	if err != nil {
		return mapPartyMemberDatabaseError(err)
	}

	return nil
}

func (r *PartyMemberRepository) Delete(
	ctx context.Context,
	member *models.PartyMember,
) error {

	return r.db.
		WithContext(ctx).
		Delete(member).
		Error()
}

func (r *PartyMemberRepository) SyncRoles(
	ctx context.Context,
	memberID uuid.UUID,
	roles []enum.PartyMemberRole,
) error {

	return r.transactionManager.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			if err := tx.
				Where(
					"party_member_id = ?",
					memberID,
				).
				Delete(
					&models.PartyMemberRole{},
				).
				Error(); err != nil {

				return err
			}

			if len(roles) == 0 {
				return nil
			}

			memberRoles := make(
				[]models.PartyMemberRole,
				0,
				len(roles),
			)

			for _, role := range roles {

				memberRoles = append(
					memberRoles,
					models.PartyMemberRole{
						ID:            uuid.New(),
						PartyMemberID: memberID,
						Role:          role,
					},
				)
			}

			if err := tx.
				Create(
					&memberRoles,
				).
				Error(); err != nil {

				return mapPartyMemberDatabaseError(err)
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

	case pgErr.Code == "23505" &&
		pgErr.ConstraintName == "idx_party_member_role":

		return appErrors.ErrInvalidPartyMemberRole

	case pgErr.Code == "23514":

		return appErrors.ErrInvalidPartyMemberRole

	default:

		return err
	}
}
