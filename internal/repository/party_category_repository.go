package repository

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyCategoryRepository struct {
	db database.DBExecutor
}

func NewPartyCategoryRepository(
	db database.DBExecutor,
) *PartyCategoryRepository {
	return &PartyCategoryRepository{
		db: db,
	}
}

func (r *PartyCategoryRepository) Replace(
	tx database.DBExecutor,
	partyID uuid.UUID,
	categoryIDs []uuid.UUID,
) error {

	var party models.Party

	if err := tx.
		Where("id = ?", partyID).
		First(&party).
		Error(); err != nil {
		return err
	}

	var categories []models.PartyCategory

	if len(categoryIDs) > 0 {
		if err := tx.
			Where("id IN ?", categoryIDs).
			Find(&categories).
			Error(); err != nil {
			return err
		}
	}

	return tx.
		Model(&party).
		Association("Categories").
		Replace(&categories)
}
