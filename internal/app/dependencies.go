package app

import (
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/service"
	auth_service "github.com/reinp/event-platform/backend/internal/service/auth"
)

func BuildDependencies(
	cfg *config.Config,
	executor database.DBExecutor,
) (*dependencies.Container, error) {

	// manager
	transactionManager :=
		database.NewTransactionManager(
			executor,
		)

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
			transactionManager,
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

	emailService := service.NewEmailService(
		clients.Mailer,
		cfg.FrontendURL,
	)

	tokenService :=
		auth_service.NewTokenService(
			userRepository,
			refreshTokenRepository,
			clients.JWT,
			clients.Clock,
			cfg.RefreshTokenDuration,
		)

	verificationService :=
		auth_service.NewVerificationService(
			userRepository,
			emailVerificationRepository,
			tokenService,
			clients.Clock,
			emailService,
			cfg.EmailVerificationDuration,
			cfg.EmailVerificationCooldown,
		)

	sessionService :=
		auth_service.NewSessionService(
			userRepository,
			refreshTokenRepository,
			tokenService,
		)

	passwordService :=
		auth_service.NewPasswordService(
			userRepository,
			passwordResetRepository,
			refreshTokenRepository,
			clients.PasswordValidator,
			emailService,
			clients.Clock,
			cfg.PasswordResetDuration,
			cfg.PasswordResetCooldown,
		)

	registrationService :=
		auth_service.NewRegistrationService(
			userRepository,
			emailVerificationRepository,
			verificationService,
			emailService,
			clients.PasswordValidator,
		)

	authService :=
		service.NewAuthService(
			userRepository,
			registrationService,
			sessionService,
			tokenService,
			verificationService,
			passwordService,
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
			transactionManager,
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
			partyRepository,
			partyMemberRepository,
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

	paypalClient := paypal.NewClient(
		cfg.PayPalClientID,
		cfg.PayPalClientSecret,
		cfg.PayPalBaseURL,
		cfg.PayPalReturnURL,
		cfg.PayPalCancelURL,
		cfg.PayPalWebhookID,
	)

	refundService := service.NewRefundService()

	paymentService :=
		service.NewPaymentService(
			purchaseService,
			ticketService,
			permissionService,
			paypalClient,
			paymentEventRepository,
			purchaseRepository,
			refundService,
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
