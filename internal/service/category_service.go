package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	partyCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/party_category"
)

type CategoryService struct {
	categories *partyCategoryRepository.Facade
}

func NewCategoryService(
	categories *partyCategoryRepository.Facade,
) *CategoryService {

	return &CategoryService{
		categories: categories,
	}
}

func (s *CategoryService) FindAll(
	ctx context.Context,
) ([]models.PartyCategory, error) {

	return s.categories.Repository.FindAll(ctx)
}

func (s *CategoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.PartyCategory, error) {

	return s.categories.Repository.FindByID(ctx, id)
}

func (s *CategoryService) Create(
	ctx context.Context,
	category *models.PartyCategory,
) error {

	return s.categories.Repository.Create(ctx, category)
}

func (s *CategoryService) Update(
	ctx context.Context,
	category *models.PartyCategory,
) error {

	return s.categories.Repository.Update(ctx, category)
}

func (s *CategoryService) Delete(
	ctx context.Context,
	category *models.PartyCategory,
) error {

	return s.categories.Repository.Delete(ctx, category)
}

func (s *CategoryService) FindPaginatedByPopularity(
	ctx context.Context,
	limit int,
) ([]models.PartyCategory, error) {

	return s.categories.Repository.FindPaginatedByPopularity(
		ctx,
		limit,
	)
}
