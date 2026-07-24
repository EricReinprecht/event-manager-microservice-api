package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/dto"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type TicketService struct {
	tickets       *repository.TicketRepository
	parties       *repository.PartyRepository
	categories    *repository.TicketCategoryRepository
	partyMembers  *repository.PartyMemberRepository
	ticketScans   *repository.TicketScanRepository
	accessWindows *repository.TicketAccessWindowRepository
	db            database.DBExecutor
}

func NewTicketService(
	ticketRepository *repository.TicketRepository,
	partyRepository *repository.PartyRepository,
	categoryRepository *repository.TicketCategoryRepository,
	partyMemberRepository *repository.PartyMemberRepository,
	ticketScansRepository *repository.TicketScanRepository,
	accessWindowRepository *repository.TicketAccessWindowRepository,
	db database.DBExecutor,
) *TicketService {

	return &TicketService{
		tickets:       ticketRepository,
		parties:       partyRepository,
		categories:    categoryRepository,
		partyMembers:  partyMemberRepository,
		ticketScans:   ticketScansRepository,
		accessWindows: accessWindowRepository,
		db:            db,
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
	scannerID uuid.UUID,
	code string,
) (*models.TicketScan, error) {

	ticket, err := s.tickets.FindByCode(
		ctx,
		code,
	)

	if err != nil {
		return nil, err
	}

	member, err := s.partyMembers.FindByPartyAndUser(
		ctx,
		ticket.TicketCategory.PartyID,
		scannerID,
	)

	if err != nil {
		return nil, appErrors.ErrNotAllowed
	}

	if member.Role != enum.RoleOrganizer &&
		member.Role != enum.RoleAdmin &&
		member.Role != enum.RoleStaff {

		return nil, appErrors.ErrNotAllowed
	}

	now := time.Now().UTC()

	window, err := s.accessWindows.FindCurrent(
		ctx,
		ticket.TicketCategoryID,
		now,
	)

	if err != nil {
		return nil, appErrors.ErrTicketNotValidNow
	}

	alreadyScanned := s.ticketScans.ExistsVerifiedInWindow(
		ctx,
		ticket.ID,
		window.StartsAt,
		window.EndsAt,
	)

	if alreadyScanned {

		return nil, appErrors.ErrTicketAlreadyScanned

	}

	status := enum.TicketScanVerified

	if ticket.TicketCategory.RequiresVerification {

		status = enum.TicketScanPending

	}

	scan := &models.TicketScan{

		TicketID: ticket.ID,

		ScannedByID: scannerID,

		ScannedAt: now,

		Status: status,
	}

	err = s.ticketScans.Create(
		ctx,
		scan,
	)

	if err != nil {
		return nil, err
	}

	return scan, nil
}

func (s *TicketService) VerifyScan(
	ctx context.Context,
	scanID uuid.UUID,
	staffID uuid.UUID,
	approved bool,
) error {

	scan, err := s.ticketScans.FindByID(
		ctx,
		scanID,
	)

	if err != nil {
		return err
	}

	if scan.Status != enum.TicketScanPending {

		return appErrors.ErrTicketScanAlreadyDecided

	}

	now := time.Now().UTC()

	if approved {

		scan.Status = enum.TicketScanVerified

	} else {

		scan.Status = enum.TicketScanRejected

	}

	scan.VerifiedAt = &now

	scan.VerifiedByID = &staffID

	return s.ticketScans.Update(
		ctx,
		scan,
	)
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
