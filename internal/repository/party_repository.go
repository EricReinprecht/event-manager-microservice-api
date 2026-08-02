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
	ctx context.Context,
	party *models.Party,
) error {

	return r.db.
		WithContext(ctx).
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
		Preload("Categories").
		Preload("Thumbnail").
		Preload("Images").
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
		Preload("Categories").
		Preload("Organizer").
		First(
			&party,
			"id = ?",
			id,
		).
		Error()

	if err != nil {

		return nil, err
	}

	return &party, nil
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
		Updates(map[string]any{

			"title": party.Title,

			"description": party.Description,

			"thumbnail_id": party.ThumbnailID,

			"location_name": party.LocationName,

			// Location metadata
			"street": party.Street,

			"house_number": party.HouseNumber,

			"city": party.City,

			"country": party.Country,

			"postal_code": party.PostalCode,

			"latitude": party.Latitude,

			"longitude": party.Longitude,

			"timezone": party.Timezone,

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
		Delete(
			party,
		).
		Error()
}
