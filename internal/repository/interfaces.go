package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketRepositoryInterface interface {
	Create(
		ctx context.Context,
		ticket *models.Ticket,
	) error

	FindByCode(
		ctx context.Context,
		code string,
	) (*models.Ticket, error)

	FindByUser(
		ctx context.Context,
		userID uuid.UUID,
	) ([]models.Ticket, error)

	CountByCategory(
		ctx context.Context,
		categoryID uuid.UUID,
	) (int64, error)

	CancelByPurchase(
		ctx context.Context,
		purchaseID uuid.UUID,
	) error
}
