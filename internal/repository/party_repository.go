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
) error {

	return tx.
		Create(party).
		Error()
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
		First(
			&party,
			"id = ?",
			id,
		).
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
		Where(
			"id = ?",
			party.ID,
		).
		Updates(map[string]interface{}{

			"title": party.Title,

			"description": party.Description,

			"location": party.Location,

			"category_id": party.CategoryID,

			"thumbnail_id": party.ThumbnailID,

			"start_at": party.StartAt,

			"end_at": party.EndAt,
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
