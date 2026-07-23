package repository

import (
	"context"

	"gorm.io/gorm"

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
