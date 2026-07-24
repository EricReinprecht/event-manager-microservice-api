package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
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
	db database.DBExecutor,
	id uuid.UUID,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := db.
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

func (r *PurchaseRepository) UpdatePayment(
	db database.DBExecutor,
	purchase *models.Purchase,
	provider string,
	paymentID string,
) error {

	purchase.PaymentProvider = provider
	purchase.PaymentID = paymentID

	return db.Save(purchase).Error()
}

func (r *PurchaseRepository) FindByPaymentID(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := r.db.
		WithContext(ctx).
		Preload("Items").
		Where(
			"payment_id = ?",
			paymentID,
		).
		First(&purchase).
		Error()

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (r *PurchaseRepository) UpdateStatus(
	db database.DBExecutor,
	purchase *models.Purchase,
	status enum.PurchaseStatus,
) error {

	purchase.Status = status

	return db.Save(purchase).Error()
}

func (r *PurchaseRepository) Find(
	ctx context.Context,
	id uuid.UUID,
) (*models.Purchase, error) {

	return r.FindByID(
		r.db.WithContext(ctx),
		id,
	)
}

func (r *PurchaseRepository) FindByPaymentIDWithDB(
	db database.DBExecutor,
	paymentID string,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := db.
		Preload("Items").
		Where(
			"payment_id = ?",
			paymentID,
		).
		First(&purchase).
		Error()

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}
