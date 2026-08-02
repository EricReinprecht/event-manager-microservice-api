package purchase_repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
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
	ctx context.Context,
	purchase *models.Purchase,
) error {

	return r.db.
		WithContext(ctx).
		Create(purchase).
		Error()
}

func (r *PurchaseRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := r.db.
		WithContext(ctx).
		Preload("Items").
		Preload("Items.TicketCategory").
		Preload("Items.TicketCategory.RefundPolicy").
		First(
			&purchase,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, mapPurchaseDatabaseError(err)
	}

	return &purchase, nil
}

func (r *PurchaseRepository) FindByIDForUpdate(
	ctx context.Context,
	id uuid.UUID,
) (*models.Purchase, error) {

	var purchase models.Purchase

	err := r.db.
		WithContext(ctx).
		Clauses(
			clause.Locking{
				Strength: "UPDATE",
			},
		).
		Preload("Items").
		Preload("Items.TicketCategory").
		Preload("Items.TicketCategory.RefundPolicy").
		First(
			&purchase,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		return nil, mapPurchaseDatabaseError(err)
	}

	return &purchase, nil
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
		First(
			&purchase,
		).
		Error()

	if err != nil {
		return nil, mapPurchaseDatabaseError(err)
	}

	return &purchase, nil
}

func (r *PurchaseRepository) FindByPaymentIDForUpdate(
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
		Preload("Items").
		Preload("Items.TicketCategory").
		Preload("Items.TicketCategory.RefundPolicy").
		Where(
			"payment_id = ?",
			paymentID,
		).
		First(
			&purchase,
		).
		Error()

	if err != nil {
		return nil, mapPurchaseDatabaseError(err)
	}

	return &purchase, nil
}

func (r *PurchaseRepository) Update(
	ctx context.Context,
	purchase *models.Purchase,
) error {

	return r.db.
		WithContext(ctx).
		Save(purchase).
		Error()
}

func (r *PurchaseRepository) UpdatePayment(
	ctx context.Context,
	purchase *models.Purchase,
	provider string,
	paymentID string,
) error {

	return r.db.
		WithContext(ctx).
		Model(purchase).
		Updates(
			map[string]any{
				"payment_provider": provider,
				"payment_id":       paymentID,
			},
		).
		Error()
}

func (r *PurchaseRepository) ReservedQuantity(
	ctx context.Context,
	categoryID uuid.UUID,
) (int64, error) {

	var quantity int64

	err := r.db.
		WithContext(ctx).
		Model(&models.PurchaseItem{}).
		Select("COALESCE(SUM(quantity),0)").
		Joins(
			"JOIN purchases ON purchases.id = purchase_items.purchase_id",
		).
		Where(
			"purchase_items.ticket_category_id = ?",
			categoryID,
		).
		Where(
			"purchases.status = ?",
			enum.PurchaseStatusPending,
		).
		Where(
			"purchases.expires_at > NOW()",
		).
		Scan(
			&quantity,
		).
		Error()

	return quantity, err
}

func mapPurchaseDatabaseError(
	err error,
) error {

	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return appErrors.ErrPurchaseNotFound
	}

	return err
}
