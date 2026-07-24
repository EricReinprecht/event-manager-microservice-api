package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketRepository struct {
	db database.DBExecutor
}

func NewTicketRepository(
	db database.DBExecutor,
) *TicketRepository {

	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) Create(
	ctx context.Context,
	ticket *models.Ticket,
) error {

	return r.db.
		WithContext(ctx).
		Create(ticket).
		Error()
}

func (r *TicketRepository) FindByCode(
	ctx context.Context,
	code string,
) (*models.Ticket, error) {

	var ticket models.Ticket

	err := r.db.
		WithContext(ctx).
		Preload(
			"TicketCategory",
			func(db *gorm.DB) *gorm.DB {
				return db.Unscoped()
			},
		).
		Preload(
			"TicketCategory.Party",
			func(db *gorm.DB) *gorm.DB {
				return db.Unscoped()
			},
		).
		Where(
			"code = ?",
			code,
		).
		First(&ticket).
		Error()

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Ticket, error) {

	var ticket models.Ticket

	err := r.db.
		WithContext(ctx).
		First(&ticket, "id = ?", id).
		Error()

	return &ticket, err
}

func (r *TicketRepository) Update(
	ctx context.Context,
	ticket *models.Ticket,
) error {

	return r.db.
		WithContext(ctx).
		Save(ticket).
		Error()
}

func (r *TicketRepository) CreatePurchase(
	ctx context.Context,
	purchase *models.Purchase,
) error {

	return r.db.
		WithContext(ctx).
		Create(purchase).
		Error()
}

func (r *TicketRepository) CreatePurchaseItem(
	ctx context.Context,
	item *models.PurchaseItem,
) error {

	return r.db.
		WithContext(ctx).
		Create(item).
		Error()
}

func (r *TicketRepository) FindByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	err := r.db.
		WithContext(ctx).
		Preload("TicketCategory").
		Preload("TicketCategory.Party").
		Preload("User").
		Find(&tickets).
		Error()

	return tickets, err
}

func (r *TicketRepository) CountByCategory(
	ctx context.Context,
	categoryID uuid.UUID,
) (int64, error) {

	var count int64

	err := r.db.
		WithContext(ctx).
		Model(&models.Ticket{}).
		Where(
			"ticket_category_id = ?",
			categoryID,
		).
		Count(&count).
		Error()

	return count, err
}
