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
) {

	authHandler := handlers.NewAuthHandler(authService)

	partyHandler := handlers.NewPartyHandler(
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
}
