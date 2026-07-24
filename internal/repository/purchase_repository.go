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
	db database.DBExecutor,
	purchase *models.Purchase,
) error {

	return db.
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
	db database.DBExecutor,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := db.
		First(
			&category,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *PurchaseRepository) Transaction(
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

		tx.Rollback()

		return err
	}

	return tx.Commit()
}
