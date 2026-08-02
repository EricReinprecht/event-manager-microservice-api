package ticket_category_repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm/clause"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketCategoryRepository struct {
	db database.DBExecutor
}

func NewTicketCategoryRepository(
	db database.DBExecutor,
) *TicketCategoryRepository {

	return &TicketCategoryRepository{
		db: db,
	}
}

func (r *TicketCategoryRepository) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	err := r.db.
		WithContext(ctx).
		Create(
			category,
		).
		Error()
	return mapDatabaseError(err)
}

func (r *TicketCategoryRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Preload("Party").
		Preload("AccessWindows").
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

func (r *TicketCategoryRepository) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	var categories []models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Preload("AccessWindows").
		Where(
			"party_id = ?",
			partyID,
		).
		Find(
			&categories,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *TicketCategoryRepository) FindByIDForUpdate(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Clauses(
			clause.Locking{
				Strength: "UPDATE",
			},
		).
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

func (r *TicketCategoryRepository) Update(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	err := r.db.
		WithContext(ctx).
		Model(
			&models.TicketCategory{},
		).
		Where(
			"id = ?",
			category.ID,
		).
		Updates(
			map[string]any{
				"name": category.Name,

				"price": category.Price,

				"capacity": category.Capacity,

				"requires_verification": category.RequiresVerification,

				"refund_requires_approval": category.RefundRequiresApproval,

				"refund_policy_id": category.RefundPolicyID,
			},
		).
		Error()
	return mapDatabaseError(err)
}

func mapDatabaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return appErrors.ErrTicketCategoryExists
	}
	return err
}

func (r *TicketCategoryRepository) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return r.db.
		WithContext(ctx).
		Delete(
			category,
		).
		Error()
}
