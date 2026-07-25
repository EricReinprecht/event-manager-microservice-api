package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Where(
			"payment_id = ?",
			paymentID,
		).
		First(&purchase).
		Error()

	if err != nil {

		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {

			return nil, appErrors.ErrPurchaseNotFound
		}

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

func (r *PurchaseRepository) FindByPaymentIDForUpdate(
	tx database.DBExecutor,
	paymentID string,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := tx.
		WithContext(context.Background()).
		Clauses(
			clause.Locking{
				Strength: "UPDATE",
			},
		).
		Preload("Items").
		Where(
			"payment_id = ?",
			paymentID,
		).
		First(
			&purchase,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (r *PurchaseRepository) ConfirmPaymentAtomic(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := r.db.
		WithContext(ctx).
		Clauses(
			clause.Locking{
				Strength: "UPDATE",
			},
		).
		Where(
			"payment_id = ?",
			paymentID,
		).
		First(
			&purchase,
		).
		Error()

	if err != nil {
		return nil, err
	}

	if purchase.Status == enum.StatusPaid {
		return &purchase, nil
	}

	purchase.Status = enum.StatusPaid

	if err := r.db.
		WithContext(ctx).
		Save(
			&purchase,
		).
		Error(); err != nil {

		return nil, err
	}

	return &purchase, nil
}

func (r *PurchaseRepository) Update(
	tx database.DBExecutor,
	purchase *models.Purchase,
) error {

	return tx.
		Save(
			purchase,
		).
		Error()
}

func (r *PurchaseRepository) ConfirmPayment(
	ctx context.Context,
	paymentID string,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := r.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			p, err := r.FindByPaymentIDForUpdate(
				tx,
				paymentID,
			)

			if err != nil {
				return err
			}

			// already processed
			if p.Status == enum.StatusPaid {
				purchase = *p
				return nil
			}

			p.Status = enum.StatusPaid

			if err := r.Update(
				tx,
				p,
			); err != nil {
				return err
			}

			purchase = *p

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}
