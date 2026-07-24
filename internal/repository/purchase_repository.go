package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PurchaseRepository struct {
	db database.DBExecutor
}

func NewPurchaseRepository(
	db database.DBExecutor,
) *PurchaseRepository {

	return &PurchaseRepository{
		db: db,
	}
}

func (r *PurchaseRepository) Create(
	ctx context.Context,
	purchase *models.Purchase,
) error {

	return r.db.WithContext(ctx).
		Create(purchase).
		Error()
}

func (r *PurchaseRepository) FindByID(
	id uuid.UUID,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := r.db.
		Preload("Items").
		First(&purchase, id).
		Error()

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (r *PurchaseRepository) FindTicketCategory(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.WithContext(ctx).
		First(&category, "id = ?", id).
		Error()

	if err != nil {
		return nil, err
	}

	return &category, nil
}
