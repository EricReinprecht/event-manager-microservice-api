package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/middleware"
)

func Register(
	router *gin.Engine,
	deps *dependencies.Container,
) {

	api := router.Group("/api")

	registerPublicRoutes(
		router,
		api,
		deps,
	)

	protected := api.Group(
		"",
		middleware.Auth(
			deps.AuthService,
		),
	)

	partyOwner := protected.Group(
		"",
		middleware.PartyOwnerMiddleware(
			deps.PermissionService,
		),
	)

	registerUserRoutes(
		protected,
		deps,
	)

	registerPartyRoutes(
		protected,
		partyOwner,
		deps,
	)

	registerCategoryRoutes(
		protected,
		deps,
	)

	registerTicketCategoryRoutes(
		protected,
		partyOwner,
		deps,
	)

	registerTicketRoutes(
		protected,
		deps,
	)

	registerPurchaseRoutes(
		protected,
		deps,
	)

	registerPaymentRoutes(
		protected,
		deps,
	)

	registerPartyMemberRoutes(
		protected,
		deps,
	)

	registerMediaRoutes(
		router,
		deps,
	)
}
