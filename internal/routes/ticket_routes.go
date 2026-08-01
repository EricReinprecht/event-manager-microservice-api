package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerTicketRoutes(
	protected *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewTicketHandler(
		deps.TicketService,
	)

	protected.GET(
		constants.MyTickets,
		handler.GetMyTickets,
	)

	protected.POST(
		constants.TicketScan,
		handler.Scan,
	)

	protected.POST(
		constants.TicketVerifyScan,
		handler.VerifyScan,
	)
}
