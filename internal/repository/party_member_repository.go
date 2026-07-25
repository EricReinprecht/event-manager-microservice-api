package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyMemberRepository struct {
	db database.DBExecutor
}

func NewPartyMemberRepository(
	db database.DBExecutor,
) *PartyMemberRepository {

	return &PartyMemberRepository{
		db: db,
	}
}

func (r *PartyMemberRepository) DB(
	ctx context.Context,
) database.DBExecutor {

	return r.db.
		WithContext(ctx)
}

func (r *PartyMemberRepository) Create(
	tx database.DBExecutor,
	member *models.PartyMember,
) error {

	return tx.
		Create(member).
		Error()
}

func (r *PartyMemberRepository) FindByPartyAndUser(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
) (*models.PartyMember, error) {

	var member models.PartyMember

	err := r.db.WithContext(ctx).
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

	err := r.db.WithContext(ctx).
		Where(
			"party_id = ?",
			partyID,
		).
		Find(&members).
		Error()

	fmt.Println("PARTY ID:", partyID)
	fmt.Println("FOUND:", len(members))
	fmt.Printf("%+v\n", members)

	return members, err
}

func (r *PartyMemberRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.PartyMember, error) {

	var member models.PartyMember

	err := r.db.WithContext(ctx).
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

func (r *PartyMemberRepository) Delete(
	ctx context.Context,
	member *models.PartyMember,
) error {

	return r.db.WithContext(ctx).
		Delete(member).
		Error()
}

func (r *PartyMemberRepository) Update(
	ctx context.Context,
	member *models.PartyMember,
) error {

	return r.db.WithContext(ctx).
		Save(member).
		Error()
}

func (r *PartyMemberRepository) Transaction(
	ctx context.Context,
	fn func(tx database.DBExecutor) error,
) error {

	tx := r.db.
		WithContext(ctx).
		Begin()

	if tx.Error() != nil {
		return tx.Error()
	}

	if err := fn(tx); err != nil {

		tx.Rollback()

		return err
	}

	return tx.Commit()
}
