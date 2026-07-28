package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/service"
)

func Register(
	router *gin.Engine,
	authLimiter *middleware.IPRateLimiter,
	authService *service.AuthService,
	userService *service.UserService,
	partyService *service.PartyService,
	categoryService *service.CategoryService,
	mediaService *service.MediaService,
	ticketCategoryService *service.TicketCategoryService,
	ticketService *service.TicketService,
	purchaseService *service.PurchaseService,
	paymentService *service.PaymentService,
	partyMemberService *service.PartyMemberService,
) {

	authHandler := handlers.NewAuthHandler(authService)

	userHandler := handlers.NewUserHandler(
		userService,
	)

	partyHandler := handlers.NewPartyHandler(
		partyService,
	)

	categoryHandler := handlers.NewCategoryHandler(
		categoryService,
	)

	mediaHandler := handlers.NewMediaHandler(
		mediaService,
	)

	ticketCategoryHandler := handlers.NewTicketCategoryHandler(
		ticketCategoryService,
		partyService,
	)

	ticketHandler := handlers.NewTicketHandler(
		ticketService,
	)

	purchaseHandler := handlers.NewPurchaseHandler(
		purchaseService,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	partyMemberHandler := handlers.NewPartyMemberHandler(
		partyMemberService,
	)

	// Public routes
	router.GET(
		"/health",
		handlers.Health,
	)

	router.POST(
		"/api/auth/register",
		authLimiter.Middleware(),
		authHandler.Register,
	)

	router.GET(
		"/api/auth/verify-email",
		authLimiter.Middleware(),
		authHandler.VerifyEmail,
	)

	router.POST(
		"/api/auth/login",
		authLimiter.Middleware(),
		authHandler.Login,
	)

	// Protected routes
	protected := router.Group("/api")

	protected.Use(
		middleware.Auth(
			authService,
		),
	)

	// * User * //
	protected.GET(
		"/users/me",
		userHandler.Me,
	)

	protected.PUT(
		"/users/me/profile",
		userHandler.CompleteProfile,
	)
	// * User * //

	// * Parties * //
	protected.POST(
		"/parties",
		partyHandler.Create,
	)

	router.GET(
		"/api/parties",
		partyHandler.GetAll,
	)

	router.GET(
		"/api/parties/:id",
		partyHandler.GetByID,
	)

	protected.PUT(
		"/parties/:id",
		partyHandler.Update,
	)

	protected.DELETE(
		"/parties/:id",
		partyHandler.Delete,
	)
	// * Parties * //

	// * Categories * //
	router.GET(
		"/api/categories",
		categoryHandler.GetAll,
	)

	router.POST(
		"/api/categories",
		categoryHandler.Create,
	)

	router.GET(
		"/api/categories/:id",
		categoryHandler.GetByID,
	)

	router.PUT(
		"/api/categories/:id",
		categoryHandler.Update,
	)

	router.DELETE(
		"/api/categories/:id",
		categoryHandler.Delete,
	)
	// * Categories * //

	// * TicketCategories * //
	protected.POST(
		"/parties/:id/ticket-categories",
		ticketCategoryHandler.Create,
	)

	router.GET(
		"/api/parties/:id/ticket-categories",
		ticketCategoryHandler.GetAll,
	)

	router.GET(
		"/api/ticket-categories/:id",
		ticketCategoryHandler.GetByID,
	)

	protected.PUT(
		"/ticket-categories/:id",
		ticketCategoryHandler.Update,
	)

	protected.DELETE(
		"/ticket-categories/:id",
		ticketCategoryHandler.Delete,
	)
	// * TicketCategories * //

	// * Tickets * //
	protected.GET(
		"/tickets/me",
		ticketHandler.GetMyTickets,
	)

	protected.POST(
		"/tickets/scan",
		ticketHandler.Scan,
	)

	protected.POST(
		"/tickets/scan/:id/verify",
		ticketHandler.VerifyScan,
	)
	// * Tickets * //

	// * Purchase * //
	protected.POST(
		"/parties/:id/purchase",
		purchaseHandler.Create,
	)

	protected.GET(
		"/purchases/:id",
		purchaseHandler.GetByID,
	)
	// * Purchase * //

	// * Payments * //
	protected.POST(
		"/purchases/:id/checkout",
		paymentHandler.CreateCheckout,
	)

	protected.POST(
		"/purchases/:id/refund",
		paymentHandler.Refund,
	)

	router.POST(
		"/api/payments/paypal/webhook",
		paymentHandler.Webhook,
	)
	// * Payments * //

	// * PartyMembers * //
	protected.POST(
		"/parties/:id/members",
		partyMemberHandler.Create,
	)

	protected.GET(
		"/parties/:id/members",
		partyMemberHandler.GetAll,
	)

	protected.DELETE(
		"/parties/:id/members/:memberID",
		partyMemberHandler.Delete,
	)

	protected.PUT(
		"/parties/:id/members/:memberID/roles",
		partyMemberHandler.UpdateRoles,
	)
	// * PartyMembers * //

	// * Media * //
	router.POST(
		"/api/media/upload",
		mediaHandler.Upload,
	)

}
