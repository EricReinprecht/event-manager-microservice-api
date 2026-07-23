package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type CategoryService struct {
	repository *repository.CategoryRepository
}

func NewCategoryService(
	repository *repository.CategoryRepository,
) *CategoryService {

	return &CategoryService{
		repository: repository,
	}
}

func (s *CategoryService) FindAll(
	ctx context.Context,
) ([]models.Category, error) {

	return s.repository.FindAll(ctx)
}

func (s *CategoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Category, error) {

	return s.repository.FindByID(ctx, id)
}

func (s *CategoryService) Create(
	ctx context.Context,
	category *models.Category,
) error {

	return s.repository.Create(ctx, category)
}

func (s *CategoryService) Update(
	ctx context.Context,
	category *models.Category,
) error {

	return s.repository.Update(ctx, category)
}

func (s *CategoryService) Delete(
	ctx context.Context,
	category *models.Category,
) error {

	return s.repository.Delete(ctx, category)
}
