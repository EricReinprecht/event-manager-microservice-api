package app

import (
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
	partyRepository "github.com/reinp/event-platform/backend/internal/repository/party"
	partyMemberRepository "github.com/reinp/event-platform/backend/internal/repository/party_member"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_category"
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

	ticketRepository :=
		repository.NewTicketRepository(
			executor,
		)

	ticketScanRepository :=
		repository.NewTicketScanRepository(
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

	ticketCategoryRepositories :=
		ticketCategoryRepository.NewFacade(
			executor,
		)

	partyMemberRepositories :=
		partyMemberRepository.NewFacade(
			executor,
			transactionManager,
		)

	partyRepositories :=
		partyRepository.NewFacade(
			executor,
			transactionManager,
			partyImageRepository,
			partyCategoryRepository,
			ticketCategoryRepositories.Write,
		)
	ticketUnitOfWork :=
		repository.NewTicketUnitOfWork(
			transactionManager,
		)

	refundUnitOfWork :=
		repository.NewRefundUnitOfWork(
			transactionManager,
		)

	purchaseUnitOfWork :=
		repository.NewPurchaseUnitOfWork(
			transactionManager,
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
			partyMemberRepositories,
			partyRepositories,
		)

	partyCRUDService := service.NewPartyCRUDService(
		partyRepositories,
		categoryRepository,
		mediaRepository,
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
			partyRepositories,
			partyMemberRepositories,
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
			ticketCategoryRepositories,
		)

	ticketService :=
		service.NewTicketService(
			ticketRepository,
			ticketScanRepository,
			permissionService,
			ticketUnitOfWork,
			clients.Clock,
			cfg.TicketVerificationTTL,
		)

	purchaseService :=
		service.NewPurchaseService(
			purchaseRepository,
			purchaseUnitOfWork,
			clients.Clock,
			cfg.PurchasePendingDuration,
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
			refundUnitOfWork,
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
