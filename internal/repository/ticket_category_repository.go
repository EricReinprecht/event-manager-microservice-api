package repository

import (
	"context"

	"github.com/reinp/event-platform/backend/internal/models"
	"gorm.io/gorm"
)

type TicketCategoryRepository struct {
	db *gorm.DB
}

func NewTicketCategoryRepository(db *gorm.DB) *TicketCategoryRepository {
	return &TicketCategoryRepository{
		db: db,
	}
}

func (r *TicketCategoryRepository) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return r.db.
		WithContext(ctx).
		Create(category).
		Error
}
