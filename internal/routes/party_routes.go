package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerPartyRoutes(
	protected *gin.RouterGroup,
	partyOwner *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewPartyHandler(
		deps.PartyService,
	)

	protected.POST(
		constants.PartyCreate,
		handler.Create,
	)

	protected.GET(
		constants.PartyList,
		handler.GetAll,
	)

	protected.GET(
		constants.PartyByID,
		handler.GetByID,
	)

	partyOwner.PUT(
		constants.PartyUpdate,
		handler.Update,
	)

	partyOwner.DELETE(
		constants.PartyDelete,
		handler.Delete,
	)

	protected.POST(
		constants.PartyPublish,
		handler.Publish,
	)

}
