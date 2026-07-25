package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketRepositoryInterface interface {
	Create(
		ctx context.Context,
		ticket *models.Ticket,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*models.Ticket, error)

	FindByCode(
		ctx context.Context,
		code string,
	) (*models.Ticket, error)

	FindByUser(
		ctx context.Context,
		userID uuid.UUID,
	) ([]models.Ticket, error)

	CancelByPurchase(
		tx database.DBExecutor,
		purchaseID uuid.UUID,
	) error
}
