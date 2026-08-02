package party_category_repository

import (
	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyCategoryWriteRepository struct{}

func NewPartyCategoryWriteRepository() *PartyCategoryWriteRepository {

	return &PartyCategoryWriteRepository{}
}

func (r *PartyCategoryWriteRepository) ReplaceForParty(
	tx database.DBExecutor,
	partyID uuid.UUID,
	categoryIDs []uuid.UUID,
) error {

	var party models.Party

	if err := tx.
		First(
			&party,
			"id = ?",
			partyID,
		).
		Error(); err != nil {

		return err
	}

	categories := make(
		[]models.PartyCategory,
		0,
		len(categoryIDs),
	)

	if len(categoryIDs) > 0 {

		if err := tx.
			Where(
				"id IN ?",
				categoryIDs,
			).
			Find(
				&categories,
			).
			Error(); err != nil {

			return err
		}

		if len(categories) != len(categoryIDs) {
			return appErrors.ErrCategoryNotFound
		}
	}

	return tx.
		Model(
			&party,
		).
		Association(
			"Categories",
		).
		Replace(
			&categories,
		)
}
