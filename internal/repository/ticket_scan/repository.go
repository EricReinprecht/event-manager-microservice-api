package ticket_scan_repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"gorm.io/gorm"
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

	err := r.db.WithContext(ctx).
		Preload("Ticket").
		Preload("Ticket.TicketCategory").
		First(
			&scan,
			"id = ?",
			id,
		).
		Error()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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

func (r *TicketScanRepository) ExistsForWindow(
	ctx context.Context,
	ticketID uuid.UUID,
	windowID uuid.UUID,
) (bool, error) {

	var count int64

	err := r.db.
		WithContext(ctx).
		Model(&models.TicketScan{}).
		Where(
			"ticket_id = ?",
			ticketID,
		).
		Where(
			"ticket_access_window_id = ?",
			windowID,
		).
		Where(
			"status IN ?",
			[]enum.TicketScanStatus{
				enum.TicketScanPending,
				enum.TicketScanVerified,
			},
		).
		Count(&count).
		Error()

	return count > 0, err
}

func (r *TicketScanRepository) FindLatestVerifiedInWindow(
	ctx context.Context,
	ticketID uuid.UUID,
	windowID *uuid.UUID,
	now time.Time,
) (*models.TicketScan, error) {

	var scan models.TicketScan

	query := r.db.
		WithContext(ctx).
		Where(
			"ticket_id = ? AND status = ? AND verified_until > ?",
			ticketID,
			enum.TicketScanVerified,
			now,
		)

	if windowID == nil {
		query = query.Where("ticket_access_window_id IS NULL")
	} else {
		query = query.Where("ticket_access_window_id = ?", *windowID)
	}

	err := query.
		First(&scan).
		Error()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &scan, nil
}

func (r *TicketScanRepository) ExistsPendingInWindow(
	ctx context.Context,
	ticketID uuid.UUID,
	windowID *uuid.UUID,
) (bool, error) {

	var count int64

	query := r.db.
		WithContext(ctx).
		Model(&models.TicketScan{}).
		Where(
			"ticket_id = ? AND status = ?",
			ticketID,
			enum.TicketScanPending,
		)

	if windowID == nil {
		query = query.Where("ticket_access_window_id IS NULL")
	} else {
		query = query.Where("ticket_access_window_id = ?", *windowID)
	}

	err := query.
		Count(&count).
		Error()

	return count > 0, err
}

func (r *TicketScanRepository) UpdateIfPending(
	ctx context.Context,
	id uuid.UUID,
	updates map[string]interface{},
) error {

	exec := r.db.
		WithContext(ctx).
		Model(&models.TicketScan{}).
		Where(
			"id = ? AND status = ?",
			id,
			enum.TicketScanPending,
		).
		Updates(updates)

	if exec.Error() != nil {
		return exec.Error()
	}

	if exec.RowsAffected() == 0 {
		return appErrors.ErrTicketScanAlreadyDecided
	}

	return nil
}
