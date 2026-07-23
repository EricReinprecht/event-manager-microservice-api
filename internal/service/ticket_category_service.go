package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

var ErrTicketCategoryExists = errors.New(
	"ticket category already exists",
)

type TicketCategoryService struct {
	ticketCategories *repository.TicketCategoryRepository
}

func NewTicketCategoryService(
	ticketCategories *repository.TicketCategoryRepository,
) *TicketCategoryService {

	return &TicketCategoryService{
		ticketCategories: ticketCategories,
	}
}

func (s *TicketCategoryService) Create(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	err := s.ticketCategories.Create(
		ctx,
		category,
	)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" {
				return ErrTicketCategoryExists
			}
		}

		return err
	}

	return nil
}

func (s *TicketCategoryService) FindByParty(
	ctx context.Context,
	partyID uuid.UUID,
) ([]models.TicketCategory, error) {

	return s.ticketCategories.FindByParty(
		ctx,
		partyID,
	)
}

func (s *TicketCategoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketCategory, error) {

	return s.ticketCategories.FindByID(
		ctx,
		id,
	)
}

func (s *TicketCategoryService) Update(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return s.ticketCategories.Update(
		ctx,
		category,
	)
}

func (s *TicketCategoryService) Delete(
	ctx context.Context,
	category *models.TicketCategory,
) error {

	return s.ticketCategories.Delete(
		ctx,
		category,
	)
}
