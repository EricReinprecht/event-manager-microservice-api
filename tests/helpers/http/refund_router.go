package http

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/handlers"
)

func RefundRouter(
	handler *handlers.PaymentHandler,
	userID interface{},
) *gin.Engine {

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			c.Set(
				"userID",
				userID,
			)

			handler.Refund(c)
		},
	)

	return router
}
