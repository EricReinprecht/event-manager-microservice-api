package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerPartyMemberRoutes(
	protected *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewPartyMemberHandler(
		deps.PartyMemberService,
	)

	protected.POST(
		constants.PartyMembers,
		handler.Create,
	)

	protected.GET(
		constants.PartyMembers,
		handler.GetAll,
	)

	protected.DELETE(
		constants.PartyMemberByID,
		handler.Delete,
	)

	protected.PUT(
		constants.PartyMemberRoles,
		handler.UpdateRoles,
	)
}
