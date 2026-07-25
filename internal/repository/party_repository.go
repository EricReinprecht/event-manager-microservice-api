package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyRepository struct {
	db database.DBExecutor
}

func NewPartyRepository(
	db database.DBExecutor,
) *PartyRepository {

	return &PartyRepository{
		db: db,
	}
}

func (r *PartyRepository) Create(
	tx database.DBExecutor,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	if err := tx.
		Create(party).
		Error(); err != nil {

		return err
	}

	for _, imageID := range imageIDs {

		link := models.PartyMedia{

			PartyID: party.ID,

			MediaID: imageID,
		}

		if err := tx.
			Create(&link).
			Error(); err != nil {

			return err
		}
	}

	return nil
}

func (r *PartyRepository) FindAll(
	ctx context.Context,
) ([]models.Party, error) {

	var parties []models.Party

	err := r.db.
		WithContext(ctx).
		Preload("Organizer").
		Preload("Category").
		Preload("Thumbnail").
		Find(&parties).
		Error()

	return parties, err
}

func (r *PartyRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Party, error) {

	var party models.Party

	err := r.db.
		WithContext(ctx).
		Preload("Thumbnail").
		Preload("Images").
		Preload("Category").
		Preload("Organizer").
		First(&party, "id = ?", id).
		Error()

	return &party, err
}

func (r *PartyRepository) Update(
	ctx context.Context,
	party *models.Party,
) error {

	return r.db.
		WithContext(ctx).
		Model(&models.Party{}).
		Where("id = ?", party.ID).
		Updates(map[string]interface{}{
			"title":        party.Title,
			"description":  party.Description,
			"location":     party.Location,
			"category_id":  party.CategoryID,
			"thumbnail_id": party.ThumbnailID,
			"start_at":     party.StartAt,
			"end_at":       party.EndAt,
		}).
		Error()
}

func (r *PartyRepository) Delete(
	ctx context.Context,
	party *models.Party,
) error {

	return r.db.
		WithContext(ctx).
		Delete(party).
		Error()
}

func (r *PartyRepository) UpdateImages(
	ctx context.Context,
	party *models.Party,
	imageIDs []uuid.UUID,
) error {

	return r.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			// remove old images

			err := tx.
				Where(
					"party_id = ?",
					party.ID,
				).
				Delete(
					&models.PartyMedia{},
				).
				Error()

			if err != nil {
				return err
			}

			// insert new relations

			for _, imageID := range imageIDs {

				link := models.PartyMedia{

					PartyID: party.ID,

					MediaID: imageID,
				}

				err := tx.
					Create(
						&link,
					).
					Error()

				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (r *PartyRepository) Transaction(
	ctx context.Context,
	fn func(tx database.DBExecutor) error,
) error {

	tx := r.db.
		WithContext(ctx).
		Begin()

	if err := tx.Error(); err != nil {
		return err
	}

	if err := fn(tx); err != nil {

		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}
