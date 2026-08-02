package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/dto"
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
	req dto.CreateCategoryRequest,
) (*models.PartyCategory, error) {
	category := &models.PartyCategory{Name: req.Name}

	if err := s.categories.Repository.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateCategoryRequest,
) (*models.PartyCategory, error) {
	category, err := s.categories.Repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	category.Name = req.Name

	if err := s.categories.Repository.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	category, err := s.categories.Repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
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
