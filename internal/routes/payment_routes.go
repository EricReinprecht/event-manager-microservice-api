package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerPaymentRoutes(
	protected *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewPaymentHandler(
		deps.PaymentService,
	)

	protected.POST(
		constants.Checkout,
		handler.CreateCheckout,
	)

	protected.POST(
		constants.Refund,
		handler.Refund,
	)

	// webhook stays public
	// because PayPal calls it without JWT

}
