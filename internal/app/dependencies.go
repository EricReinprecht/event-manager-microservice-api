package app

import (
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
)

func BuildDependencies(
	cfg *config.Config,
	executor database.DBExecutor,
) (*dependencies.Container, error) {

	// repositories

	userRepository :=
		repository.NewUserRepository(
			executor,
		)

	emailVerificationRepository :=
		repository.NewEmailVerificationRepository(
			executor,
		)

	refreshTokenRepository :=
		repository.NewRefreshTokenRepository(
			executor,
		)

	passwordResetRepository :=
		repository.NewPasswordResetTokenRepository(
			executor,
		)

	partyRepository :=
		repository.NewPartyRepository(
			executor,
		)

	partyQueryRepository :=
		repository.NewPartyQueryRepository(
			executor,
		)

	partyImageRepository :=
		repository.NewPartyImageRepository(
			executor,
		)

	categoryRepository :=
		repository.NewCategoryRepository(
			executor,
		)

	partyCategoryRepository :=
		repository.NewPartyCategoryRepository(
			executor,
		)

	mediaRepository :=
		repository.NewMediaRepository(
			executor,
		)

	partyMemberRepository :=
		repository.NewPartyMemberRepository(
			executor,
		)

	partyMemberRoleRepository :=
		repository.NewPartyMemberRoleRepository(
			executor,
		)

	ticketCategoryRepository :=
		repository.NewTicketCategoryRepository(
			executor,
		)

	ticketRepository :=
		repository.NewTicketRepository(
			executor,
		)

	ticketScanRepository :=
		repository.NewTicketScanRepository(
			executor,
		)

	ticketAccessWindowRepository :=
		repository.NewTicketAccessWindowRepository(
			executor,
		)

	purchaseRepository :=
		repository.NewPurchaseRepository(
			executor,
		)

	paymentEventRepository :=
		repository.NewPaymentEventRepository(
			executor,
		)

	// clients

	clients :=
		NewClients(
			cfg,
		)

	// services

	emailService :=
		service.NewEmailService(
			clients.Mailer,
			cfg.FrontendURL,
		)

	authService :=
		service.NewAuthService(
			userRepository,
			emailVerificationRepository,
			refreshTokenRepository,
			passwordResetRepository,
			clients.JWT,
			emailService,
			clients.PasswordValidator,
			cfg.RefreshTokenDuration,
			cfg.PasswordResetCooldown,
		)

	userService :=
		service.NewUserService(
			userRepository,
			refreshTokenRepository,
			clients.PasswordValidator,
			passwordResetRepository,
			emailVerificationRepository,
		)

	partyMemberService :=
		service.NewPartyMemberService(
			partyMemberRepository,
			partyRepository,
			partyMemberRoleRepository,
		)

	partyCRUDService :=
		service.NewPartyCRUDService(
			partyRepository,
			partyImageRepository,
			partyMemberRepository,
			partyMemberRoleRepository,
			categoryRepository,
			partyCategoryRepository,
			mediaRepository,
			ticketCategoryRepository,
			database.NewTransactionManager(
				executor,
			),
		)

	partyQueryService :=
		service.NewPartyQueryService(
			partyQueryRepository,
		)

	partyAccessService :=
		service.NewPartyAccessService(
			partyQueryRepository,
		)

	partyService :=
		service.NewPartyService(
			partyCRUDService,
			partyQueryService,
			partyAccessService,
		)

	permissionService :=
		service.NewPermissionService(
			partyService,
			partyMemberService,
		)

	categoryService :=
		service.NewCategoryService(
			categoryRepository,
		)

	mediaService :=
		service.NewMediaService(
			mediaRepository,
		)

	ticketCategoryService :=
		service.NewTicketCategoryService(
			ticketCategoryRepository,
		)

	ticketService :=
		service.NewTicketService(
			ticketRepository,
			partyMemberRepository,
			ticketScanRepository,
			ticketAccessWindowRepository,
			executor,
			clients.Clock,
			cfg.TicketVerificationTTL,
		)

	purchaseService :=
		service.NewPurchaseService(
			purchaseRepository,
			ticketRepository,
		)

	paymentService :=
		service.NewPaymentService(
			purchaseService,
			ticketService,
			partyMemberService,
			clients.PayPalClient,
			paymentEventRepository,
			purchaseRepository,
			service.NewRefundService(),
		)

	return &dependencies.Container{

		AuthService: authService,

		UserService: userService,

		PartyService: partyService,

		CategoryService: categoryService,

		MediaService: mediaService,

		TicketCategoryService: ticketCategoryService,

		TicketService: ticketService,

		PurchaseService: purchaseService,

		PaymentService: paymentService,

		PartyMemberService: partyMemberService,

		PermissionService: permissionService,

		RefreshTokenDuration: cfg.RefreshTokenDuration,

		CookieSecure: cfg.CookieSecure,
	}, nil
}
