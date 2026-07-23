package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

var ErrTicketAlreadyUsed = errors.New(
	"ticket already used",
)

type TicketService struct {
	tickets    *repository.TicketRepository
	parties    *repository.PartyRepository
	categories *repository.TicketCategoryRepository
	db         database.DBExecutor
}

func NewTicketService(
	ticketRepository *repository.TicketRepository,
	partyRepository *repository.PartyRepository,
	categoryRepository *repository.TicketCategoryRepository,
	db database.DBExecutor,
) *TicketService {

	return &TicketService{
		tickets:    ticketRepository,
		parties:    partyRepository,
		categories: categoryRepository,
		db:         db,
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

	tx := s.db.WithContext(ctx).Begin()

	ticketRepo := repository.NewTicketRepository(tx)
	categoryRepo := repository.NewTicketCategoryRepository(tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	_, err := s.parties.FindByID(ctx, partyID)

	if err != nil {
		tx.Rollback()
		return nil, appErrors.ErrPartyNotFound
	}

	purchase := &models.Purchase{
		UserID:  userID,
		PartyID: partyID,
		Status:  enum.StatusPending,
	}

	err = ticketRepo.CreatePurchase(ctx, purchase)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range items {

		category, err := categoryRepo.FindByIDForUpdate(
			ctx,
			item.TicketCategoryID,
		)

		if err != nil {
			tx.Rollback()
			return nil, appErrors.ErrTicketCategoryNotFound
		}

		createdTickets, err := ticketRepo.CountByCategory(
			ctx,
			category.ID,
		)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		available := int64(category.Capacity) - createdTickets

		if available < int64(item.Quantity) {

			tx.Rollback()

			return nil, appErrors.ErrNotEnoughTickets
		}

		err = ticketRepo.CreatePurchaseItem(
			ctx,
			&models.PurchaseItem{
				PurchaseID:       purchase.ID,
				TicketCategoryID: category.ID,
				Quantity:         item.Quantity,
				Price:            category.Price,
			},
		)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		for i := 0; i < item.Quantity; i++ {

			ticket := &models.Ticket{
				Code: uuid.NewString(),

				TicketCategoryID: category.ID,

				UserID: userID,
			}

			err = ticketRepo.Create(
				ctx,
				ticket,
			)

			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
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
