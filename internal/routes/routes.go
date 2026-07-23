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
	// router.GET(
	// 	"/api/categories",
	// 	ticketCategoryHandler.GetAll,
	// )

	protected.POST(
		"/parties/:id/ticket-categories",
		ticketCategoryHandler.Create,
	)

	// router.GET(
	// 	"/api/categories/:id",
	// 	ticketCategoryHandler.GetByID,
	// )

	// router.PUT(
	// 	"/api/categories/:id",
	// 	ticketCategoryHandler.Update,
	// )

	// router.DELETE(
	// 	"/api/categories/:id",
	// 	ticketCategoryHandler.Delete,
	// )
	// * TicketCategories * //

	// * Media * //
	router.POST(
		"/api/media/upload",
		mediaHandler.Upload,
	)

}
