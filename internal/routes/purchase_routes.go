package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerPurchaseRoutes(
	protected *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewPurchaseHandler(
		deps.PurchaseService,
	)

	protected.POST(
		constants.PurchaseCreate,
		handler.Create,
	)

	protected.GET(
		constants.PurchaseByID,
		handler.GetByID,
	)
}
