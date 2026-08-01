package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerTicketCategoryRoutes(
	protected *gin.RouterGroup,
	partyOwner *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewTicketCategoryHandler(
		deps.TicketCategoryService,
		deps.PermissionService,
	)

	partyOwner.POST(
		constants.PartyTicketCategories,
		handler.Create,
	)

	protected.GET(
		constants.PartyTicketCategories,
		handler.GetAll,
	)

	protected.GET(
		constants.TicketCategoryByID,
		handler.GetByID,
	)

	partyOwner.PUT(
		constants.TicketCategoryUpdate,
		handler.Update,
	)

	partyOwner.DELETE(
		constants.TicketCategoryUpdate,
		handler.Delete,
	)

}
