package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(
	db *gorm.DB,
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
		Error
}

func (r *MediaRepository) FindByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.Media, error) {

	var media []models.Media

	err := r.db.
		WithContext(ctx).
		Where("id IN ?", ids).
		Find(&media).
		Error

	return media, err
}
