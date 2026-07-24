package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketAccessWindowRepository struct {
	db database.DBExecutor
}

func NewTicketAccessWindowRepository(
	db database.DBExecutor,
) *TicketAccessWindowRepository {

	return &TicketAccessWindowRepository{
		db: db,
	}
}

func (r *TicketAccessWindowRepository) IsAllowedNow(
	ctx context.Context,
	categoryID uuid.UUID,
	now time.Time,
) bool {

	var count int64

	r.db.
		WithContext(ctx).
		Model(&models.TicketAccessWindow{}).
		Where(
			"ticket_category_id = ? AND starts_at <= ? AND ends_at >= ?",
			categoryID,
			now,
			now,
		).
		Count(&count)

	return count > 0
}

func (r *TicketAccessWindowRepository) FindCurrent(
	ctx context.Context,
	categoryID uuid.UUID,
	now time.Time,
) (*models.TicketAccessWindow, error) {

	now = now.UTC()

	var window models.TicketAccessWindow

	err := r.db.
		WithContext(ctx).
		Where(
			"ticket_category_id = ? AND starts_at <= ? AND ends_at >= ?",
			categoryID,
			now,
			now,
		).
		First(&window).
		Error()

	if err != nil {
		return nil, err
	}

	return &window, nil
}
