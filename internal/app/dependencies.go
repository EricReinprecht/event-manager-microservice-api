package app

import (
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/repository"
	emailVerificationRepository "github.com/reinp/event-platform/backend/internal/repository/email_verification"
	mediaRepository "github.com/reinp/event-platform/backend/internal/repository/media"
	partyRepository "github.com/reinp/event-platform/backend/internal/repository/party"
	partyCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/party_category"
	partyImageRepository "github.com/reinp/event-platform/backend/internal/repository/party_image"
	partyMemberRepository "github.com/reinp/event-platform/backend/internal/repository/party_member"
	passwordResetTokenRepository "github.com/reinp/event-platform/backend/internal/repository/password_reset_token"
	paymentEventRepository "github.com/reinp/event-platform/backend/internal/repository/payment_event"
	purchaseRepository "github.com/reinp/event-platform/backend/internal/repository/purchase"
	refreshTokenRepository "github.com/reinp/event-platform/backend/internal/repository/refresh_token"
	ticketRepository "github.com/reinp/event-platform/backend/internal/repository/ticket"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_category"
	ticketScanRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_scan"
	userRepository "github.com/reinp/event-platform/backend/internal/repository/user"
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
	userRepositories :=
		userRepository.NewFacade(
			executor,
		)
	userRepository := userRepositories.Repository

	emailVerificationRepositories :=
		emailVerificationRepository.NewFacade(
			executor,
		)
	emailVerificationRepository := emailVerificationRepositories.Repository

	refreshTokenRepositories :=
		refreshTokenRepository.NewFacade(
			executor,
		)
	refreshTokenRepository := refreshTokenRepositories.Repository

	passwordResetRepositories :=
		passwordResetTokenRepository.NewFacade(
			executor,
		)
	passwordResetRepository := passwordResetRepositories.Repository

	mediaRepositories :=
		mediaRepository.NewFacade(
			executor,
		)
	mediaRepository := mediaRepositories.Repository

	partyImageRepositories :=
		partyImageRepository.NewFacade(
			executor,
		)
	partyImageRepository := partyImageRepositories.Repository

	ticketRepositories :=
		ticketRepository.NewFacade(
			executor,
		)
	ticketRepository := ticketRepositories.Repository

	ticketScanRepositories :=
		ticketScanRepository.NewFacade(
			executor,
		)
	ticketScanRepository := ticketScanRepositories.Repository

	paymentEventRepositories :=
		paymentEventRepository.NewFacade(
			executor,
		)
	paymentEventRepository := paymentEventRepositories.Repository

	ticketCategoryRepositories :=
		ticketCategoryRepository.NewFacade(
			executor,
		)

	partyMemberRepositories :=
		partyMemberRepository.NewFacade(
			executor,
			transactionManager,
		)

	partyCategoryRepositories :=
		partyCategoryRepository.NewFacade(
			executor,
		)

	partyRepositories :=
		partyRepository.NewFacade(
			executor,
			transactionManager,
			partyImageRepository,
			partyCategoryRepositories.Write,
			ticketCategoryRepositories.Write,
		)

	purchaseRepositories :=
		purchaseRepository.NewFacade(
			executor,
			transactionManager,
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
		partyCategoryRepositories,
		mediaRepository,
	)

	partyQueryService :=
		service.NewPartyQueryService(
			partyRepositories,
		)

	partyAccessService :=
		service.NewPartyAccessService(
			partyRepositories,
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
			partyCategoryRepositories,
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
			purchaseRepositories,
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
