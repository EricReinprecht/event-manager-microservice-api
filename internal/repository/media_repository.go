package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type MediaRepository struct {
	db database.DBExecutor
}

func NewMediaRepository(
	db database.DBExecutor,
) *MediaRepository {

	return &MediaRepository{
		db: db,
	}
}

func (r *MediaRepository) Create(
	ctx context.Context,
	media *models.Media,
) error {

	return r.db.
		WithContext(ctx).
		Create(media).
		Error()
}

func (r *MediaRepository) FindByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Media, error) {

	var media []models.Media

	err := r.db.
		WithContext(ctx).
		Where(
			"id IN ?",
			ids,
		).
		Find(&media).
		Error()

	return media, err
}

func (r *MediaRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Media, error) {

	var media models.Media

	err := r.db.
		WithContext(ctx).
		First(
			&media,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &media, nil
}
