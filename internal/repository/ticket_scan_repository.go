package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketScanRepository struct {
	db database.DBExecutor
}

func NewTicketScanRepository(
	db database.DBExecutor,
) *TicketScanRepository {

	return &TicketScanRepository{
		db: db,
	}
}

func (r *TicketScanRepository) Create(
	ctx context.Context,
	scan *models.TicketScan,
) error {

	return r.db.
		WithContext(ctx).
		Create(scan).
		Error()
}

func (r *TicketScanRepository) ExistsInWindow(
	ctx context.Context,
	ticketID uuid.UUID,
	start time.Time,
	end time.Time,
) bool {

	var count int64

	r.db.
		WithContext(ctx).
		Model(&models.TicketScan{}).
		Where(
			"ticket_id = ? AND scanned_at >= ? AND scanned_at <= ?",
			ticketID,
			start,
			end,
		).
		Count(&count)

	return count > 0
}
