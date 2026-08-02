package party_category_repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type PartyCategoryRepository struct {
	db database.DBExecutor
}

func NewPartyCategoryRepository(
	db database.DBExecutor,
) *PartyCategoryRepository {

	return &PartyCategoryRepository{
		db: db,
	}
}

func (r *PartyCategoryRepository) FindAll(
	ctx context.Context,
) ([]models.PartyCategory, error) {

	var categories []models.PartyCategory

	err := r.db.
		WithContext(ctx).
		Find(
			&categories,
		).
		Error()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PartyCategoryRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.PartyCategory, error) {

	var category models.PartyCategory

	err := r.db.
		WithContext(ctx).
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

func (r *PartyCategoryRepository) Create(
	ctx context.Context,
	category *models.PartyCategory,
) error {

	return r.db.
		WithContext(ctx).
		Create(
			category,
		).
		Error()
}

func (r *PartyCategoryRepository) Update(
	ctx context.Context,
	category *models.PartyCategory,
) error {

	return r.db.
		WithContext(ctx).
		Model(
			&models.PartyCategory{},
		).
		Where(
			"id = ?",
			category.ID,
		).
		Updates(
			map[string]any{
				"name": category.Name,
			},
		).
		Error()
}

func (r *PartyCategoryRepository) Delete(
	ctx context.Context,
	category *models.PartyCategory,
) error {

	return r.db.
		WithContext(ctx).
		Delete(
			category,
		).
		Error()
}

func (r *PartyCategoryRepository) FindPaginatedByPopularity(
	ctx context.Context,
	limit int,
) ([]models.PartyCategory, error) {

	var categories []models.PartyCategory

	err := r.db.
		WithContext(ctx).
		Model(
			&models.PartyCategory{},
		).
		Select(
			"categories.*, COUNT(DISTINCT parties.id) AS usage_count",
		).
		Joins(
			`
			LEFT JOIN party_categories
				ON party_categories.category_id = categories.id
			`,
		).
		Joins(
			`
			LEFT JOIN parties
				ON parties.id = party_categories.party_id
				AND parties.deleted_at IS NULL
			`,
		).
		Group(
			"categories.id",
		).
		Order(
			"usage_count DESC",
		).
		Limit(
			limit,
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
