package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func Register(
	router *gin.Engine,
	deps *dependencies.Container,
) {

	authHandler := handlers.NewAuthHandler(
		deps.AuthService,
	)

	userHandler := handlers.NewUserHandler(
		deps.UserService,
	)

	sessionHandler := handlers.NewSessionHandler(
		deps.AuthService,
	)

	partyHandler := handlers.NewPartyHandler(
		deps.PartyService,
	)

	categoryHandler := handlers.NewCategoryHandler(
		deps.CategoryService,
	)

	mediaHandler := handlers.NewMediaHandler(
		deps.MediaService,
	)

	ticketCategoryHandler := handlers.NewTicketCategoryHandler(
		deps.TicketCategoryService,
		deps.PartyService,
	)

	ticketHandler := handlers.NewTicketHandler(
		deps.TicketService,
	)

	purchaseHandler := handlers.NewPurchaseHandler(
		deps.PurchaseService,
	)

	paymentHandler := handlers.NewPaymentHandler(
		deps.PaymentService,
	)

	partyMemberHandler := handlers.NewPartyMemberHandler(
		deps.PartyMemberService,
	)

	authLimiter := middleware.NewIPRateLimiter(
		1.0/12,
		5,
	)

	refreshLimiter := middleware.NewIPRateLimiter(
		0.5,
		20,
	)

	// =========================
	// Public
	// =========================

	router.GET(
		"/health",
		handlers.Health,
	)

	api := router.Group(constants.API)

	api.POST(
		constants.AuthRegister,
		authLimiter.Middleware(),
		authHandler.Register,
	)

	api.GET(
		constants.AuthVerifyEmail,
		authLimiter.Middleware(),
		authHandler.VerifyEmail,
	)

	api.POST(
		constants.AuthLogin,
		authLimiter.Middleware(),
		authHandler.Login,
	)

	api.POST(
		constants.AuthRefresh,
		refreshLimiter.Middleware(),
		authHandler.Refresh,
	)

	api.POST(
		constants.AuthLogout,
		authLimiter.Middleware(),
		authHandler.Logout,
	)

	api.POST(
		constants.AuthForgotPassword,
		refreshLimiter.Middleware(),
		authHandler.ForgotPassword,
	)

	api.POST(
		constants.AuthResetPassword,
		authLimiter.Middleware(),
		authHandler.ResetPassword,
	)

	api.POST(
		constants.AuthResendVerification,
		authLimiter.Middleware(),
		authHandler.ResendVerificationEmail,
	)

	// =========================
	// Protected
	// =========================

	protected := api.Use(
		middleware.Auth(
			deps.AuthService,
		),
	)

	// Users

	protected.GET(
		constants.UserSessions,
		sessionHandler.GetSessions,
	)

	protected.DELETE(
		constants.UserSessionByFamilyID,
		sessionHandler.DeleteSession,
	)

	protected.DELETE(
		constants.UserSessions,
		sessionHandler.LogoutAll,
	)

	protected.GET(
		constants.UserMe,
		userHandler.Me,
	)

	protected.PUT(
		constants.UserCompleteProfile,
		userHandler.CompleteProfile,
	)

	protected.PUT(
		constants.UserPassword,
		userHandler.ChangePassword,
	)

	protected.GET(
		constants.UserParties,
		partyHandler.GetMyParties,
	)

	// Parties

	protected.POST(
		constants.PartyCreate,
		partyHandler.Create,
	)

	protected.GET(
		constants.PartyList,
		partyHandler.GetAll,
	)

	protected.GET(
		constants.PartyByID,
		partyHandler.GetByID,
	)

	protected.PUT(
		constants.PartyUpdate,
		partyHandler.Update,
	)

	protected.DELETE(
		constants.PartyDelete,
		partyHandler.Delete,
	)

	// Categories

	protected.GET(
		constants.CategoryList,
		categoryHandler.GetAll,
	)

	protected.GET(
		constants.CategoryListPopular,
		categoryHandler.GetPaginatedByPopularity,
	)

	protected.POST(
		constants.CategoryCreate,
		categoryHandler.Create,
	)

	protected.GET(
		constants.CategoryByID,
		categoryHandler.GetByID,
	)

	protected.PUT(
		constants.CategoryByID,
		categoryHandler.Update,
	)

	protected.DELETE(
		constants.CategoryByID,
		categoryHandler.Delete,
	)

	// Ticket Categories

	protected.POST(
		constants.PartyTicketCategories,
		ticketCategoryHandler.Create,
	)

	router.GET(
		constants.PartyTicketCategories,
		ticketCategoryHandler.GetAll,
	)

	router.GET(
		constants.TicketCategoryByID,
		ticketCategoryHandler.GetByID,
	)

	protected.PUT(
		constants.TicketCategoryUpdate,
		ticketCategoryHandler.Update,
	)

	protected.DELETE(
		constants.TicketCategoryUpdate,
		ticketCategoryHandler.Delete,
	)

	// Tickets

	protected.GET(
		constants.MyTickets,
		ticketHandler.GetMyTickets,
	)

	protected.POST(
		constants.TicketScan,
		ticketHandler.Scan,
	)

	protected.POST(
		constants.TicketVerifyScan,
		ticketHandler.VerifyScan,
	)

	// Purchases

	protected.POST(
		constants.PurchaseCreate,
		purchaseHandler.Create,
	)

	protected.GET(
		constants.PurchaseByID,
		purchaseHandler.GetByID,
	)

	// Payments

	protected.POST(
		constants.Checkout,
		paymentHandler.CreateCheckout,
	)

	protected.POST(
		constants.Refund,
		paymentHandler.Refund,
	)

	router.POST(
		constants.PayPalWebhook,
		paymentHandler.Webhook,
	)

	// Party Members

	protected.POST(
		constants.PartyMembers,
		partyMemberHandler.Create,
	)

	protected.GET(
		constants.PartyMembers,
		partyMemberHandler.GetAll,
	)

	protected.DELETE(
		constants.PartyMemberByID,
		partyMemberHandler.Delete,
	)

	protected.PUT(
		constants.PartyMemberRoles,
		partyMemberHandler.UpdateRoles,
	)

	// Media

	router.POST(
		constants.MediaUpload,
		mediaHandler.Upload,
	)
}
