package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

var ErrTicketCategoryExists = errors.New(
	"ticket category already exists",
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

	err := s.repo.Create(
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
