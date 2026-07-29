package app

import (
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type Repositories struct {

	// User / Auth
	UserRepository *repository.UserRepository

	RefreshTokenRepository *repository.RefreshTokenRepository

	EmailVerificationRepository *repository.EmailVerificationRepository

	PasswordResetTokenRepository *repository.PasswordResetTokenRepository

	// Party
	PartyRepository *repository.PartyRepository

	PartyQueryRepository *repository.PartyQueryRepository

	PartyImageRepository *repository.PartyImageRepository

	PartyMemberRepository *repository.PartyMemberRepository

	PartyMemberRoleRepository *repository.PartyMemberRoleRepository

	// Media / Categories
	CategoryRepository *repository.CategoryRepository

	MediaRepository *repository.MediaRepository

	// Ticketing
	TicketCategoryRepository *repository.TicketCategoryRepository

	TicketRepository *repository.TicketRepository

	TicketScanRepository *repository.TicketScanRepository

	TicketAccessWindowRepository *repository.TicketAccessWindowRepository

	// Purchase / Payment
	PurchaseRepository *repository.PurchaseRepository

	PaymentEventRepository *repository.PaymentEventRepository
}

func NewRepositories(
	executor database.DBExecutor,
) *Repositories {

	return &Repositories{

		// =====================
		// AUTH
		// =====================

		UserRepository: repository.NewUserRepository(
			executor,
		),

		RefreshTokenRepository: repository.NewRefreshTokenRepository(
			executor,
		),

		EmailVerificationRepository: repository.NewEmailVerificationRepository(
			executor,
		),

		PasswordResetTokenRepository: repository.NewPasswordResetTokenRepository(
			executor,
		),

		// =====================
		// PARTY
		// =====================

		PartyRepository: repository.NewPartyRepository(
			executor,
		),

		PartyQueryRepository: repository.NewPartyQueryRepository(
			executor,
		),

		PartyImageRepository: repository.NewPartyImageRepository(
			executor,
		),

		PartyMemberRepository: repository.NewPartyMemberRepository(
			executor,
		),

		PartyMemberRoleRepository: repository.NewPartyMemberRoleRepository(
			executor,
		),

		// =====================
		// MEDIA / CATEGORY
		// =====================

		CategoryRepository: repository.NewCategoryRepository(
			executor,
		),

		MediaRepository: repository.NewMediaRepository(
			executor,
		),

		// =====================
		// TICKETS
		// =====================

		TicketCategoryRepository: repository.NewTicketCategoryRepository(
			executor,
		),

		TicketRepository: repository.NewTicketRepository(
			executor,
		),

		TicketScanRepository: repository.NewTicketScanRepository(
			executor,
		),

		TicketAccessWindowRepository: repository.NewTicketAccessWindowRepository(
			executor,
		),

		// =====================
		// PURCHASE / PAYMENT
		// =====================

		PurchaseRepository: repository.NewPurchaseRepository(
			executor,
		),

		PaymentEventRepository: repository.NewPaymentEventRepository(
			executor,
		),
	}
}
