package repository

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type PartyMemberRoleRepository struct {
	db database.DBExecutor
}

func NewPartyMemberRoleRepository(
	db database.DBExecutor,
) *PartyMemberRoleRepository {

	return &PartyMemberRoleRepository{
		db: db,
	}
}

func (r *PartyMemberRoleRepository) Create(
	tx database.DBExecutor,
	role *models.PartyMemberRole,
) error {

	return tx.
		Create(role).
		Error()
}

func (r *PartyMemberRoleRepository) Delete(
	tx database.DBExecutor,
	memberID uuid.UUID,
	role enum.PartyMemberRole,
) error {

	return tx.
		Where(
			"party_member_id = ? AND role = ?",
			memberID,
			role,
		).
		Delete(
			&models.PartyMemberRole{},
		).
		Error()
}

func (r *PartyMemberRoleRepository) DeleteAll(
	tx database.DBExecutor,
	memberID uuid.UUID,
) error {

	return tx.
		Where(
			"party_member_id = ?",
			memberID,
		).
		Delete(
			&models.PartyMemberRole{},
		).
		Error()
}
