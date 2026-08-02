package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type TicketService struct {
	tickets         repository.TicketRepositoryInterface
	scans           *repository.TicketScanRepository
	permissions     *PermissionService
	unitOfWork      *repository.TicketUnitOfWork
	clock           clock.Clock
	verificationTTL time.Duration
}

func NewTicketService(
	tickets repository.TicketRepositoryInterface,
	scans *repository.TicketScanRepository,
	permissions *PermissionService,
	unitOfWork *repository.TicketUnitOfWork,
	clock clock.Clock,
	verificationTTL time.Duration,
) *TicketService {

	return &TicketService{
		tickets:         tickets,
		scans:           scans,
		permissions:     permissions,
		unitOfWork:      unitOfWork,
		clock:           clock,
		verificationTTL: verificationTTL,
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

	var createdScan *models.TicketScan

	err := s.unitOfWork.Transaction(
		ctx,
		func(repositories *repository.TicketTransactionRepositories) error {

			ticket, err := repositories.Tickets.FindByCode(
				ctx,
				code,
			)

			if err != nil {
				return err
			}

			if ticket.TicketCategory.DeletedAt.Valid {
				return appErrors.ErrTicketCategoryNotFound
			}

			if ticket.Status == enum.TicketStatusCancelled {
				return appErrors.ErrTicketNotValidNow
			}

			if err := s.permissions.RequireScanTickets(
				ctx,
				ticket.TicketCategory.PartyID,
				scannerID,
			); err != nil {

				return err
			}

			now := s.clock.Now()

			var window *models.TicketAccessWindow

			if len(ticket.TicketCategory.AccessWindows) == 0 {
				party := ticket.TicketCategory.Party
				if now.Before(party.StartAt) || now.After(party.EndAt) {
					return appErrors.ErrTicketNotValidNow
				}
			} else {
				window, err = repositories.AccessWindows.FindCurrent(
					ctx,
					ticket.TicketCategoryID,
					now,
				)

				if err != nil {
					return appErrors.ErrTicketNotValidNow
				}
			}

			var windowID *uuid.UUID
			if window != nil {
				windowID = &window.ID
			}

			if err := ensureTicketNotAlreadyScanned(
				ctx,
				repositories.Scans,
				ticket.ID,
				windowID,
				now,
			); err != nil {

				return err
			}

			scan := s.newTicketScan(
				ticket,
				window,
				scannerID,
				now,
			)

			if err := repositories.Scans.Create(
				ctx,
				scan,
			); err != nil {

				return err
			}

			createdScan = scan

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return createdScan, nil
}

func (s *TicketService) VerifyScan(
	ctx context.Context,
	scanID uuid.UUID,
	staffID uuid.UUID,
	approved bool,
) error {

	scan, err := s.scans.FindByID(
		ctx,
		scanID,
	)

	if err != nil {
		return err
	}

	if err := s.permissions.RequireScanTickets(
		ctx,
		scan.Ticket.TicketCategory.PartyID,
		staffID,
	); err != nil {

		return err
	}

	if scan.Status != enum.TicketScanPending {
		return appErrors.ErrTicketScanAlreadyDecided
	}

	now := s.clock.Now()

	if scan.VerificationExpiresAt != nil &&
		now.After(*scan.VerificationExpiresAt) {

		return appErrors.ErrTicketVerificationExpired
	}

	updates := s.scanDecisionUpdates(
		staffID,
		approved,
		now,
	)

	return s.scans.UpdateIfPending(
		ctx,
		scan.ID,
		updates,
	)
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

func (s *TicketService) GenerateFromPurchase(
	ctx context.Context,
	purchase *models.Purchase,
) error {

	tickets := makeTicketsFromPurchase(
		purchase,
	)

	return s.unitOfWork.Transaction(
		ctx,
		func(repositories *repository.TicketTransactionRepositories) error {

			for index := range tickets {

				if err := repositories.Tickets.Create(
					ctx,
					&tickets[index],
				); err != nil {

					return err
				}
			}

			return nil
		},
	)
}

func (s *TicketService) newTicketScan(
	ticket *models.Ticket,
	window *models.TicketAccessWindow,
	scannerID uuid.UUID,
	now time.Time,
) *models.TicketScan {

	scan := &models.TicketScan{
		ID:          uuid.New(),
		TicketID:    ticket.ID,
		ScannedByID: scannerID,
		ScannedAt:   now,
		Status:      enum.TicketScanPending,
	}

	if window != nil {
		scan.TicketAccessWindowID = &window.ID
	}

	expiry := now.Add(
		s.verificationTTL,
	)

	if ticket.TicketCategory.RequiresVerification {

		scan.VerificationExpiresAt = &expiry

		return scan
	}

	scan.Status = enum.TicketScanVerified
	scan.VerifiedAt = &now
	scan.VerifiedByID = &scannerID
	scan.VerifiedUntil = &expiry

	return scan
}

func (s *TicketService) scanDecisionUpdates(
	staffID uuid.UUID,
	approved bool,
	now time.Time,
) map[string]any {

	if approved {

		verifiedUntil := now.Add(
			s.verificationTTL,
		)

		return map[string]any{
			"status":                  enum.TicketScanVerified,
			"verified_at":             now,
			"verified_by_id":          staffID,
			"verified_until":          verifiedUntil,
			"verification_expires_at": nil,
		}
	}

	return map[string]any{
		"status":                  enum.TicketScanRejected,
		"verified_at":             now,
		"verified_by_id":          staffID,
		"verification_expires_at": nil,
		"verified_until":          nil,
	}
}

func ensureTicketNotAlreadyScanned(
	ctx context.Context,
	scans *repository.TicketScanRepository,
	ticketID uuid.UUID,
	windowID *uuid.UUID,
	now time.Time,
) error {

	existingVerified, err :=
		scans.FindLatestVerifiedInWindow(
			ctx,
			ticketID,
			windowID,
			now,
		)

	if err == nil && existingVerified != nil {
		return appErrors.ErrTicketAlreadyScanned
	}

	if err != nil {
		return err
	}

	pendingExists, err :=
		scans.ExistsPendingInWindow(
			ctx,
			ticketID,
			windowID,
		)

	if err != nil {
		return err
	}

	if pendingExists {
		return appErrors.ErrTicketAlreadyScanned
	}

	return nil
}

func makeTicketsFromPurchase(
	purchase *models.Purchase,
) []models.Ticket {

	totalQuantity := 0

	for _, item := range purchase.Items {
		totalQuantity += item.Quantity
	}

	tickets := make(
		[]models.Ticket,
		0,
		totalQuantity,
	)

	for _, item := range purchase.Items {

		for range item.Quantity {

			tickets = append(
				tickets,
				models.Ticket{
					ID: uuid.New(),

					Code: strings.ToUpper(
						uuid.NewString()[:8],
					),

					UserID: purchase.UserID,

					PurchaseID: purchase.ID,

					TicketCategoryID: item.TicketCategoryID,

					Status: enum.TicketStatusActive,
				},
			)
		}
	}

	return tickets
}
