package app

import (
	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/mail"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/security"
	"github.com/reinp/event-platform/backend/internal/service"
)

type Services struct {
	AuthService *service.AuthService

	UserService *service.UserService

	PartyService *service.PartyService

	CategoryService *service.CategoryService

	MediaService *service.MediaService

	TicketCategoryService *service.TicketCategoryService

	TicketService *service.TicketService

	PurchaseService *service.PurchaseService

	PaymentService *service.PaymentService

	PartyMemberService *service.PartyMemberService
}

func NewServices(
	cfg *config.Config,
	repos *Repositories,
	executor database.DBExecutor,
) *Services {

	appClock := clock.RealClock{}

	// =====================
	// EXTERNAL SERVICES
	// =====================

	jwt := auth.NewJWT(
		cfg.JWTSecret,
		appClock,
	)

	mailer := mail.NewMailer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	emailService := service.NewEmailService(
		mailer,
		cfg.FrontendURL,
	)

	passwordValidator :=
		security.NewPasswordValidator()

	// =====================
	// AUTH
	// =====================

	authService := service.NewAuthService(

		repos.UserRepository,

		repos.EmailVerificationRepository,

		repos.RefreshTokenRepository,

		repos.PasswordResetTokenRepository,

		jwt,

		emailService,

		passwordValidator,

		cfg.RefreshTokenDuration,

		cfg.PasswordResetCooldown,
	)

	userService := service.NewUserService(

		repos.UserRepository,

		repos.RefreshTokenRepository,

		passwordValidator,

		repos.PasswordResetTokenRepository,

		repos.EmailVerificationRepository,
	)

	// =====================
	// PARTY DOMAIN
	// =====================

	transactionManager :=
		database.NewTransactionManager(
			executor,
		)

	partyCRUDService :=
		service.NewPartyCRUDService(

			repos.PartyRepository,

			repos.PartyImageRepository,

			repos.PartyMemberRepository,

			repos.PartyMemberRoleRepository,

			repos.CategoryRepository,

			repos.PartyCategoryRepository,

			repos.MediaRepository,

			repos.TicketCategoryRepository,

			transactionManager,
		)

	partyQueryService :=
		service.NewPartyQueryService(

			repos.PartyQueryRepository,
		)

	partyAccessService :=
		service.NewPartyAccessService(
			repos.PartyQueryRepository,
		)

	partyService :=
		service.NewPartyService(

			partyCRUDService,

			partyQueryService,

			partyAccessService,
		)

	partyMemberService :=
		service.NewPartyMemberService(

			repos.PartyMemberRepository,

			repos.PartyRepository,

			repos.PartyMemberRoleRepository,
		)

	categoryService :=
		service.NewCategoryService(
			repos.CategoryRepository,
		)

	mediaService :=
		service.NewMediaService(
			repos.MediaRepository,
		)

	// =====================
	// TICKETS
	// =====================

	ticketCategoryService :=
		service.NewTicketCategoryService(
			repos.TicketCategoryRepository,
		)

	ticketService :=
		service.NewTicketService(

			repos.TicketRepository,

			repos.PartyMemberRepository,

			repos.TicketScanRepository,

			repos.TicketAccessWindowRepository,

			executor,

			appClock,

			cfg.TicketVerificationTTL,
		)

	// =====================
	// PURCHASE
	// =====================

	purchaseService :=
		service.NewPurchaseService(

			repos.PurchaseRepository,

			repos.TicketRepository,
		)

	refundService :=
		service.NewRefundService()

	paypalClient :=
		paypal.NewClient(

			cfg.PayPalClientID,

			cfg.PayPalClientSecret,

			cfg.PayPalBaseURL,

			cfg.PayPalReturnURL,

			cfg.PayPalCancelURL,

			cfg.PayPalWebhookID,
		)

	paymentService :=
		service.NewPaymentService(

			purchaseService,

			ticketService,

			partyMemberService,

			paypalClient,

			repos.PaymentEventRepository,

			repos.PurchaseRepository,

			refundService,
		)

	return &Services{

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
	}
}
