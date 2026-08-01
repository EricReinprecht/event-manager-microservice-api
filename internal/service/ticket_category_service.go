package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type TicketCategoryService struct {
	repo *repository.TicketCategoryRepository
}

func NewTicketCategoryService(
	repo *repository.TicketCategoryRepository,
) *TicketCategoryService {
	return &TicketCategoryService{
		repo: repo,
	}
}

func (s *TicketCategoryService) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	if err := s.validateCategory(category); err != nil {
		return err
	}

	err := s.repo.Create(ctx, category)

	if err != nil {
		return mapTicketCategoryError(err)
	}

	return nil
}

func (s *TicketCategoryService) Update(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	if err := s.validateCategory(category); err != nil {
		return err
	}

	err := s.repo.Update(ctx, category)

	if err != nil {
		return mapTicketCategoryError(err)
	}

	return nil
}

func (s *TicketCategoryService) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	return s.repo.FindByParty(
		ctx,
		partyID,
	)
}

func (s *TicketCategoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	return s.repo.FindByID(
		ctx,
		id,
	)
}

func (s *TicketCategoryService) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return s.repo.Delete(
		ctx,
		category,
	)
}

func mapTicketCategoryError(err error) error {

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case "23505":
			return appErrors.ErrTicketCategoryExists
		}
	}

	return err
}

func (s *TicketCategoryService) validateCategory(
	category *models.TicketCategory,
) error {

	if len(category.AccessWindows) == 0 {
		return appErrors.ErrTicketAccessWindowRequired
	}

	for _, window := range category.AccessWindows {

		if window.EndsAt.Before(window.StartsAt) {
			return appErrors.ErrAccessWindowInvalid
		}
	}

	return nil
}
