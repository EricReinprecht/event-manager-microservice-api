package repository

import (
	"context"
	"errors"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"gorm.io/gorm"
)

type PaymentEventRepository struct {
	db database.DBExecutor
}

func NewPaymentEventRepository(
	db database.DBExecutor,
) *PaymentEventRepository {

	return &PaymentEventRepository{
		db: db,
	}
}

func (r *PaymentEventRepository) Create(
	event *models.PaymentEvent,
) error {

	return r.db.Create(event).Error()
}

func (r *PaymentEventRepository) Update(
	event *models.PaymentEvent,
) error {

	return r.db.
		Save(event).
		Error()
}

func (r *PaymentEventRepository) FindByEventID(
	ctx context.Context,
	eventID string,
) (*models.PaymentEvent, error) {

	var event models.PaymentEvent

	result := r.db.
		WithContext(ctx).
		Where(
			"event_id = ?",
			eventID,
		).
		First(&event)

	if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error() != nil {
		return nil, result.Error()
	}

	return &event, nil
}
