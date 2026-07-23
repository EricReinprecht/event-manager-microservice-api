package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

var ErrTicketAlreadyUsed = errors.New(
	"ticket already used",
)

type TicketService struct {
	tickets *repository.TicketRepository
}

func NewTicketService(
	tickets *repository.TicketRepository,
) *TicketService {

	return &TicketService{
		tickets: tickets,
	}
}

func (s *TicketService) Create(
	ctx context.Context,
	ticket *models.Ticket,
) error {

	return s.tickets.Create(
		ctx,
		ticket,
	)
}

func (s *TicketService) FindByCode(
	ctx context.Context,
	code string,
) (*models.Ticket, error) {

	return s.tickets.FindByCode(
		ctx,
		code,
	)
}

func (s *TicketService) Scan(
	ctx context.Context,
	code string,
) (*models.Ticket, error) {

	ticket, err := s.tickets.FindByCode(
		ctx,
		code,
	)

	if err != nil {
		return nil, err
	}

	if ticket.UsedAt != nil {
		return nil, ErrTicketAlreadyUsed
	}

	return ticket, nil
}

func (s *TicketService) CreatePurchase(
	ctx context.Context,
	partyID uuid.UUID,
	userID uuid.UUID,
	items []dto.PurchaseTicketItem,
) (*models.Purchase, error) {

	purchase := &models.Purchase{
		UserID:  userID,
		PartyID: partyID,
		Status:  "pending",
	}

	err := s.tickets.CreatePurchase(
		ctx,
		purchase,
	)

	if err != nil {
		return nil, err
	}

	for _, item := range items {

		err := s.tickets.CreatePurchaseItem(
			ctx,
			&models.PurchaseItem{
				PurchaseID:       purchase.ID,
				TicketCategoryID: item.TicketCategoryID,
				Quantity:         item.Quantity,
			},
		)

		if err != nil {
			return nil, err
		}
	}

	return purchase, nil
}

func (s *TicketService) GetMyTickets(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.Ticket, error) {

	return s.tickets.FindByUser(
		ctx,
		userID,
	)
}
