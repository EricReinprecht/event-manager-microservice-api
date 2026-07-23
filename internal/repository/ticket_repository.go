package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(
	db *gorm.DB,
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
		Error
}

func (r *TicketRepository) FindByCode(
	ctx context.Context,
	code string,
) (*models.Ticket, error) {

	var ticket models.Ticket

	err := r.db.
		WithContext(ctx).
		Preload("TicketCategory").
		Preload("TicketCategory.Party").
		Preload("User").
		Where("code = ?", code).
		First(&ticket).
		Error

	return &ticket, err
}

func (r *TicketRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Ticket, error) {

	var ticket models.Ticket

	err := r.db.
		WithContext(ctx).
		First(&ticket, "id = ?", id).
		Error

	return &ticket, err
}

func (r *TicketRepository) Update(
	ctx context.Context,
	ticket *models.Ticket,
) error {

	return r.db.
		WithContext(ctx).
		Save(ticket).
		Error
}

func (r *TicketRepository) CreatePurchase(
	ctx context.Context,
	purchase *models.Purchase,
) error {

	return r.db.
		WithContext(ctx).
		Create(purchase).
		Error
}

func (r *TicketRepository) CreatePurchaseItem(
	ctx context.Context,
	item *models.PurchaseItem,
) error {

	return r.db.
		WithContext(ctx).
		Create(item).
		Error
}

func (r *TicketRepository) FindByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("TicketCategory").
		Preload("Party").
		Find(&tickets).
		Error

	return tickets, err
}
