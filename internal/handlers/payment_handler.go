package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/payment/paypal"
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

	headers := paypal.WebhookHeaders{

		TransmissionID: c.GetHeader(
			"PAYPAL-TRANSMISSION-ID",
		),

		TransmissionTime: c.GetHeader(
			"PAYPAL-TRANSMISSION-TIME",
		),

		TransmissionSig: c.GetHeader(
			"PAYPAL-TRANSMISSION-SIG",
		),

		CertURL: c.GetHeader(
			"PAYPAL-CERT-URL",
		),

		AuthAlgo: c.GetHeader(
			"PAYPAL-AUTH-ALGO",
		),
	}

	body, err := io.ReadAll(
		c.Request.Body,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "cannot read body",
			},
		)

		return
	}

	err = h.service.VerifyWebhook(
		c.Request.Context(),
		headers,
		body,
	)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid paypal signature",
			},
		)

		return
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(body, &payload); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid json",
			},
		)

		return
	}

	resource, okResource := payload["resource"].(map[string]interface{})

	if !okResource {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid webhook payload",
			},
		)

		return
	}

	paymentID, okPaymentID := resource["id"].(string)

	if !okPaymentID {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "missing payment id",
			},
		)

		return
	}

	_, err = h.service.ConfirmPayment(
		c.Request.Context(),
		paymentID,
	)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "payment confirmed",
		},
	)
}
