package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
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

func (r *TicketScanRepository) ExistsVerifiedInWindow(
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
			`
			ticket_id = ?
			AND scanned_at >= ?
			AND scanned_at <= ?
			AND status = ?
			`,
			ticketID,
			start,
			end,
			enum.TicketScanVerified,
		).
		Count(&count)

	return count > 0
}

func (r *TicketScanRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.TicketScan, error) {

	var scan models.TicketScan

	err := r.db.
		WithContext(ctx).
		First(&scan, "id = ?", id).
		Error()

	if err != nil {
		return nil, err
	}

	return &scan, nil
}

func (r *TicketScanRepository) Update(
	ctx context.Context,
	scan *models.TicketScan,
) error {

	return r.db.
		WithContext(ctx).
		Save(scan).
		Error()
}
