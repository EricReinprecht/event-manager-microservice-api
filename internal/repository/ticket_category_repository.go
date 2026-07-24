package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"gorm.io/gorm"
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

	return r.db.
		WithContext(ctx).
		Create(category).
		Error()
}

func (r *TicketCategoryRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Preload("Party").
		First(&category, "id = ?", id).
		Error()

	return &category, err
}

func (r *TicketCategoryRepository) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	var categories []models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Where("party_id = ?", partyID).
		Find(&categories).
		Error()

	return categories, err
}

func (r *TicketCategoryRepository) Update(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return r.db.
		WithContext(ctx).
		Save(category).
		Error()
}

func (r *TicketCategoryRepository) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return r.db.
		WithContext(ctx).
		Delete(category).
		Error()
}

func (r *TicketCategoryRepository) FindByIDForUpdate(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	var category models.TicketCategory

	err := r.db.
		WithContext(ctx).
		Raw(
			`
			SELECT *
			FROM ticket_categories
			WHERE id = ?
			AND deleted_at IS NULL
			FOR UPDATE
			`,
			id,
		).
		Scan(&category).
		Error()

	if err != nil {
		return nil, err
	}

	if category.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &category, nil
}
