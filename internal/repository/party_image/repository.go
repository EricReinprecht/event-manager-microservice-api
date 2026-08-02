package party_image_repository

import (
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyImageRepository struct {
	db database.DBExecutor
}

func NewPartyImageRepository(
	db database.DBExecutor,
) *PartyImageRepository {

	return &PartyImageRepository{
		db: db,
	}
}

func (r *PartyImageRepository) Replace(
	tx database.DBExecutor,
	partyID uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	if err := tx.
		Where(
			"party_id = ?",
			partyID,
		).
		Delete(
			&models.PartyMedia{},
		).
		Error(); err != nil {

		return err
	}

	for _, imageID := range imageIDs {

		link := models.PartyMedia{
			PartyID: partyID,
			MediaID: imageID,
		}

		if err := tx.
			Create(
				&link,
			).
			Error(); err != nil {

			return err
		}
	}

	return nil
}
