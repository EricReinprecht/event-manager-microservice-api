package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/middleware"
	"github.com/reinp/event-platform/backend/internal/service"
)

func Register(
	router *gin.Engine,
	authService *service.AuthService,
	partyService *service.PartyService,
	categoryService *service.CategoryService,
	mediaService *service.MediaService,
	ticketCategoryService *service.TicketCategoryService,
	ticketService *service.TicketService,
	PurchaseService *service.PurchaseService,
	partyMemberService *service.PartyMemberService,
) {

	authHandler := handlers.NewAuthHandler(authService)

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
		authHandler.Register,
	)

	router.POST(
		"/api/auth/login",
		authHandler.Login,
	)

	// Protected routes
	protected := router.Group("/api")

	protected.Use(
		middleware.Auth(
			authService.Secret(),
		),
	)

	protected.GET(
		"/users/me",
		handlers.Me,
	)

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
	protected.POST(
		"/parties/:id/tickets/purchase",
		ticketHandler.Purchase,
	)

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

	// * PartyMembers * //
	protected.POST(
		"/parties/:id/members",
		partyMemberHandler.Create,
	)

	protected.GET(
		"/parties/:id/members",
		partyMemberHandler.GetAll,
	)

	protected.PUT(
		"/parties/:id/members/:memberID",
		partyMemberHandler.Update,
	)

	protected.DELETE(
		"/parties/:id/members/:memberID",
		partyMemberHandler.Delete,
	)
	// * PartyMembers * //

	// * Media * //
	router.POST(
		"/api/media/upload",
		mediaHandler.Upload,
	)

}
