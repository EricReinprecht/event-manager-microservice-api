package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	emailVerificationRepository "github.com/reinp/event-platform/backend/internal/repository/email_verification"
	mediaRepository "github.com/reinp/event-platform/backend/internal/repository/media"
	partyImageRepository "github.com/reinp/event-platform/backend/internal/repository/party_image"
	partyMediaRepository "github.com/reinp/event-platform/backend/internal/repository/party_media"
	passwordResetTokenRepository "github.com/reinp/event-platform/backend/internal/repository/password_reset_token"
	paymentEventRepository "github.com/reinp/event-platform/backend/internal/repository/payment_event"
	refreshTokenRepository "github.com/reinp/event-platform/backend/internal/repository/refresh_token"
	ticketRepository "github.com/reinp/event-platform/backend/internal/repository/ticket"
	ticketAccessWindowRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_access_window"
	ticketScanRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_scan"
	userRepository "github.com/reinp/event-platform/backend/internal/repository/user"
)

// Deprecated compatibility aliases. New dependency wiring should use the
// entity-specific repository packages and their facades directly.
type EmailVerificationRepository = emailVerificationRepository.EmailVerificationRepository
type MediaRepository = mediaRepository.MediaRepository
type PartyImageRepository = partyImageRepository.PartyImageRepository
type PartyMediaRepository = partyMediaRepository.PartyMediaRepository
type PasswordResetTokenRepository = passwordResetTokenRepository.PasswordResetTokenRepository
type PaymentEventRepository = paymentEventRepository.PaymentEventRepository
type RefreshTokenRepository = refreshTokenRepository.RefreshTokenRepository
type TicketRepository = ticketRepository.TicketRepository
type TicketAccessWindowRepository = ticketAccessWindowRepository.TicketAccessWindowRepository
type TicketScanRepository = ticketScanRepository.TicketScanRepository
type UserRepository = userRepository.UserRepository

func NewEmailVerificationRepository(db database.DBExecutor) *EmailVerificationRepository {
	return emailVerificationRepository.NewEmailVerificationRepository(db)
}

func NewMediaRepository(db database.DBExecutor) *MediaRepository {
	return mediaRepository.NewMediaRepository(db)
}

func NewPartyImageRepository(db database.DBExecutor) *PartyImageRepository {
	return partyImageRepository.NewPartyImageRepository(db)
}

func NewPartyMediaRepository(db database.DBExecutor) *PartyMediaRepository {
	return partyMediaRepository.NewPartyMediaRepository(db)
}

func NewPasswordResetTokenRepository(db database.DBExecutor) *PasswordResetTokenRepository {
	return passwordResetTokenRepository.NewPasswordResetTokenRepository(db)
}

func NewPaymentEventRepository(db database.DBExecutor) *PaymentEventRepository {
	return paymentEventRepository.NewPaymentEventRepository(db)
}

func NewRefreshTokenRepository(db database.DBExecutor) *RefreshTokenRepository {
	return refreshTokenRepository.NewRefreshTokenRepository(db)
}

func NewTicketRepository(db database.DBExecutor) *TicketRepository {
	return ticketRepository.NewTicketRepository(db)
}

func NewTicketAccessWindowRepository(db database.DBExecutor) *TicketAccessWindowRepository {
	return ticketAccessWindowRepository.NewTicketAccessWindowRepository(db)
}

func NewTicketScanRepository(db database.DBExecutor) *TicketScanRepository {
	return ticketScanRepository.NewTicketScanRepository(db)
}

func NewUserRepository(db database.DBExecutor) *UserRepository {
	return userRepository.NewUserRepository(db)
}

type TicketRepositoryInterface interface {
	Create(
		ctx context.Context,
		ticket *models.Ticket,
	) error

	FindByCode(
		ctx context.Context,
		code string,
	) (*models.Ticket, error)

	FindByUser(
		ctx context.Context,
		userID uuid.UUID,
	) ([]models.Ticket, error)

	CountByCategory(
		ctx context.Context,
		categoryID uuid.UUID,
	) (int64, error)

	CancelByPurchase(
		ctx context.Context,
		purchaseID uuid.UUID,
	) error
}
