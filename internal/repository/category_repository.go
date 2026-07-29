package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type CategoryRepository struct {
	db database.DBExecutor
}

func NewCategoryRepository(
	db database.DBExecutor,
) *CategoryRepository {

	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) FindAll(
	ctx context.Context,
) ([]models.Category, error) {

	var categories []models.Category

	err := r.db.
		WithContext(ctx).
		Find(&categories).
		Error()

	return categories, err
}

func (r *CategoryRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Category, error) {

	var category models.Category

	err := r.db.
		WithContext(ctx).
		First(
			&category,
			"id = ?",
			id,
		).
		Error()

	return &category, err
}

func (r *CategoryRepository) Create(
	ctx context.Context,
	category *models.Category,
) error {

	return r.db.
		WithContext(ctx).
		Create(category).
		Error()
}

func (r *CategoryRepository) Update(
	ctx context.Context,
	category *models.Category,
) error {

	return r.db.
		WithContext(ctx).
		Save(category).
		Error()
}

func (r *CategoryRepository) Delete(
	ctx context.Context,
	category *models.Category,
) error {

	return r.db.
		WithContext(ctx).
		Delete(category).
		Error()
}

func (r *CategoryRepository) FindPaginatedByPopularity(
	ctx context.Context,
	limit int,
) ([]models.Category, error) {

	var categories []models.Category

	err := r.db.
		WithContext(ctx).
		Model(&models.Category{}).
		Select(
			"categories.*, COUNT(parties.id) AS usage_count",
		).
		Joins(
			"LEFT JOIN party_categories ON party_categories.category_id = categories.id",
		).
		Joins(
			"LEFT JOIN parties ON parties.id = party_categories.party_id",
		).
		Group(
			"categories.id",
		).
		Order(
			"usage_count DESC",
		).
		Limit(limit).
		Find(&categories).
		Error()

	return categories, err
}
