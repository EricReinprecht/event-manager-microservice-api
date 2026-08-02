package party_media_repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyMediaRepository struct {
	db database.DBExecutor
}

func NewPartyMediaRepository(
	db database.DBExecutor,
) *PartyMediaRepository {

	return &PartyMediaRepository{
		db: db,
	}
}

func (r *PartyMediaRepository) ReplaceImages(
	ctx context.Context,
	partyID uuid.UUID,
	imageIDs []uuid.UUID,
) error {

	tx := r.db.
		WithContext(ctx).
		Begin()

	if err := tx.Error(); err != nil {
		return err
	}

	err := tx.
		Where(
			"party_id = ?",
			partyID,
		).
		Delete(
			&models.PartyMedia{},
		).
		Error()

	if err != nil {

		tx.Rollback()

		return err
	}

	for _, imageID := range imageIDs {

		link := models.PartyMedia{

			PartyID: partyID,

			MediaID: imageID,
		}

		if err := tx.
			Create(&link).
			Error(); err != nil {

			tx.Rollback()

			return err
		}
	}

	return tx.Commit()
}
