package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(
	service *service.PaymentService,
) *PaymentHandler {

	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) CreateCheckout(c *gin.Context) {

	purchaseID, err := uuid.Parse(
		c.Param("id"),
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid purchase id",
			},
		)

		return
	}

	url, err := h.service.CreateCheckout(
		c.Request.Context(),
		purchaseID,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"approval_url": url,
		},
	)
}

func (h *PaymentHandler) Webhook(c *gin.Context) {

	var payload struct {
		PaymentID string `json:"payment_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	purchase, err := h.service.ConfirmPayment(
		c.Request.Context(),
		payload.PaymentID,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		purchase,
	)
}
