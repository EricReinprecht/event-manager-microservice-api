package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"

	"github.com/reinp/event-platform/backend/internal/models/enum"
)

type TicketService struct {
	tickets         *repository.TicketRepository
	partyMembers    *repository.PartyMemberRepository
	ticketScans     *repository.TicketScanRepository
	accessWindows   *repository.TicketAccessWindowRepository
	db              database.DBExecutor
	clock           clock.Clock
	verificationTTL time.Duration
}

func NewTicketService(
	ticketRepository *repository.TicketRepository,
	partyMemberRepository *repository.PartyMemberRepository,
	ticketScansRepository *repository.TicketScanRepository,
	accessWindowRepository *repository.TicketAccessWindowRepository,
	db database.DBExecutor,
	clock clock.Clock,
	verificationTTL time.Duration,
) *TicketService {

	return &TicketService{
		tickets:         ticketRepository,
		partyMembers:    partyMemberRepository,
		ticketScans:     ticketScansRepository,
		accessWindows:   accessWindowRepository,
		db:              db,
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

	tx := s.db.Begin()

	if tx.Error() != nil {
		return nil, tx.Error()
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	ticketRepo := repository.NewTicketRepository(tx)
	scanRepo := repository.NewTicketScanRepository(tx)
	memberRepo := repository.NewPartyMemberRepository(tx)
	windowRepo := repository.NewTicketAccessWindowRepository(tx)

	ticket, err := ticketRepo.FindByCode(
		ctx,
		code,
	)

	if err != nil {

		tx.Rollback()

		return nil, err
	}

	if ticket.TicketCategory.DeletedAt.Valid {

		tx.Rollback()

		return nil, gorm.ErrRecordNotFound
	}

	if ticket.Status == enum.TicketStatusCancelled {

		tx.Rollback()

		return nil, appErrors.ErrTicketNotValidNow
	}

	member, err := memberRepo.FindByPartyAndUser(
		ctx,
		ticket.TicketCategory.PartyID,
		scannerID,
	)

	if err != nil {
		tx.Rollback()
		return nil, appErrors.ErrNotAllowed
	}

	if member.Role != enum.RoleOrganizer &&
		member.Role != enum.RoleAdmin &&
		member.Role != enum.RoleStaff {

		tx.Rollback()
		return nil, appErrors.ErrNotAllowed
	}

	now := s.clock.Now()

	window, err := windowRepo.FindCurrent(
		ctx,
		ticket.TicketCategoryID,
		now,
	)

	if err != nil {
		tx.Rollback()
		return nil, appErrors.ErrTicketNotValidNow
	}

	// Check already verified scan in this access window
	existingVerified, err := scanRepo.FindLatestVerifiedInWindow(
		ctx,
		ticket.ID,
		window.ID,
		now,
	)

	if err == nil && existingVerified != nil {

		tx.Rollback()

		return nil, appErrors.ErrTicketAlreadyScanned
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		tx.Rollback()

		return nil, err
	}

	// Check existing pending scan in this access window
	exists, err := scanRepo.ExistsPendingInWindow(
		ctx,
		ticket.ID,
		window.ID,
	)

	if err != nil {

		tx.Rollback()

		return nil, err
	}

	if exists {

		tx.Rollback()

		return nil, appErrors.ErrTicketAlreadyScanned
	}

	status := enum.TicketScanPending

	var verifiedAt *time.Time
	var verifiedByID *uuid.UUID
	var verificationExpiresAt *time.Time
	var verifiedUntil *time.Time

	if !ticket.TicketCategory.RequiresVerification {

		status = enum.TicketScanVerified

		verifiedAt = &now
		verifiedByID = &scannerID

		expiry := now.Add(
			s.verificationTTL,
		)

		verifiedUntil = &expiry
	}

	if ticket.TicketCategory.RequiresVerification {

		expiry := now.Add(
			s.verificationTTL,
		)

		verificationExpiresAt = &expiry
	}

	scan := &models.TicketScan{

		TicketID: ticket.ID,

		TicketAccessWindowID: window.ID,

		ScannedByID: scannerID,

		ScannedAt: now,

		Status: status,

		VerifiedAt: verifiedAt,

		VerifiedByID: verifiedByID,

		VerificationExpiresAt: verificationExpiresAt,

		VerifiedUntil: verifiedUntil,
	}

	err = scanRepo.Create(
		ctx,
		scan,
	)

	if err != nil {

		tx.Rollback()

		return nil, err
	}

	if err := tx.Commit(); err != nil {

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

	// Permission check
	member, err := s.partyMembers.FindByPartyAndUser(
		ctx,
		scan.Ticket.TicketCategory.PartyID,
		staffID,
	)

	if err != nil {
		return appErrors.ErrNotAllowed
	}

	if member.Role != enum.RoleOrganizer &&
		member.Role != enum.RoleAdmin &&
		member.Role != enum.RoleStaff {

		return appErrors.ErrNotAllowed
	}

	// Only pending scans can be decided
	if scan.Status != enum.TicketScanPending {

		return appErrors.ErrTicketScanAlreadyDecided
	}

	now := s.clock.Now()

	// Pending verification expired
	if scan.VerificationExpiresAt != nil &&
		now.After(*scan.VerificationExpiresAt) {

		return appErrors.ErrTicketVerificationExpired
	}

	updates := map[string]interface{}{}

	if approved {

		verifiedUntil := now.Add(
			s.verificationTTL,
		)

		updates = map[string]interface{}{

			"status": enum.TicketScanVerified,

			"verified_at": now,

			"verified_by_id": staffID,

			"verified_until": verifiedUntil,
		}

	} else {

		updates = map[string]interface{}{

			"status": enum.TicketScanRejected,

			// decision metadata
			"verified_at": now,

			"verified_by_id": staffID,

			// clear verification validity
			"verification_expires_at": nil,

			"verified_until": nil,
		}
	}

	err = s.ticketScans.UpdateIfPending(
		ctx,
		scan.ID,
		updates,
	)

	if err != nil {

		return err
	}

	return nil
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

	tx := s.db.WithContext(ctx).Begin()

	if tx.Error() != nil {
		return tx.Error()
	}

	ticketRepo := repository.NewTicketRepository(tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range purchase.Items {

		for i := 0; i < item.Quantity; i++ {

			ticket := models.Ticket{
				ID: uuid.New(),

				Code: strings.ToUpper(
					uuid.NewString()[:8],
				),

				UserID: purchase.UserID,

				TicketCategoryID: item.TicketCategoryID,
			}

			if err := ticketRepo.Create(
				ctx,
				&ticket,
			); err != nil {

				tx.Rollback()

				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
