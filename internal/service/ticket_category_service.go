package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_category"
)

type TicketCategoryService struct {
	ticketCategories *ticketCategoryRepository.Facade
}

func NewTicketCategoryService(
	ticketCategories *ticketCategoryRepository.Facade,
) *TicketCategoryService {

	return &TicketCategoryService{
		ticketCategories: ticketCategories,
	}
}

func (s *TicketCategoryService) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	if err := s.validateCategory(category); err != nil {
		return err
	}

	err := s.ticketCategories.Repository.Create(ctx, category)

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

	err := s.ticketCategories.Repository.Update(ctx, category)

	if err != nil {
		return mapTicketCategoryError(err)
	}

	return nil
}

func (s *TicketCategoryService) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	return s.ticketCategories.Repository.FindByParty(
		ctx,
		partyID,
	)
}

func (s *TicketCategoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	return s.ticketCategories.Repository.FindByID(
		ctx,
		id,
	)
}

func (s *TicketCategoryService) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return s.ticketCategories.Repository.Delete(
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
