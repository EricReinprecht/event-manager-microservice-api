package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/models"
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

	var payload struct {
		EventID string `json:"id"`

		EventType string `json:"event_type"`

		Resource struct {
			ID string `json:"id"`

			SupplementaryData struct {
				RelatedIDs struct {
					OrderID string `json:"order_id"`
				} `json:"related_ids"`
			} `json:"supplementary_data"`
		} `json:"resource"`
	}

	err = json.Unmarshal(
		body,
		&payload,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid json",
			},
		)

		return
	}

	// ---------------------------------------
	// Idempotency
	// ---------------------------------------

	if h.processWebhookEvent(
		c,
		payload.EventID,
		payload.EventType,
		body,
	) {
		return
	}

	// ---------------------------------------
	// Capture approved orders
	// ---------------------------------------

	if payload.EventType == "CHECKOUT.ORDER.APPROVED" {

		err := h.service.CapturePayment(
			c.Request.Context(),
			payload.Resource.ID,
		)

		if err != nil {

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "capture failed",
				},
			)

			return
		}

		err = h.service.MarkPaymentEventProcessed(
			c.Request.Context(),
			payload.EventID,
		)

		if err != nil {

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "could not mark event processed",
				},
			)

			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "capture requested",
			},
		)

		return
	}

	// ---------------------------------------
	// Ignore everything except completed
	// ---------------------------------------

	if payload.EventType != "PAYMENT.CAPTURE.COMPLETED" {

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "event ignored",
			},
		)

		return
	}

	orderID := payload.
		Resource.
		SupplementaryData.
		RelatedIDs.
		OrderID

	if orderID == "" {

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "missing order id",
			},
		)

		return
	}

	// ---------------------------------------
	// Confirm payment
	// ---------------------------------------

	_, err = h.service.ConfirmPayment(
		c.Request.Context(),
		orderID,
	)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "payment processing failed",
			},
		)

		return
	}

	// ---------------------------------------
	// Mark webhook processed
	// ---------------------------------------

	err = h.service.MarkPaymentEventProcessed(
		c.Request.Context(),
		payload.EventID,
	)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "could not mark event processed",
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

func (h *PaymentHandler) processWebhookEvent(
	c *gin.Context,
	payloadID string,
	eventType string,
	body []byte,
) bool {

	existing, err := h.service.FindPaymentEvent(
		c.Request.Context(),
		payloadID,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "could not check webhook event",
			},
		)

		return true
	}

	if existing != nil && existing.Processed {

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "event already processed",
			},
		)

		return true
	}

	if existing == nil {

		event := &models.PaymentEvent{

			ID: uuid.New(),

			Provider: "paypal",

			EventID: payloadID,

			Type: eventType,

			Payload: string(body),

			Processed: false,
		}

		err = h.service.CreatePaymentEvent(
			c.Request.Context(),
			event,
		)

		if err != nil {

			c.JSON(
				http.StatusOK,
				gin.H{
					"message": "duplicate event",
				},
			)

			return true
		}
	}

	return false
}
