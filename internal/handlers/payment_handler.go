package handlers

import (
	"encoding/json"
	"io"
	"log"
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

	log.Println("========== PAYPAL WEBHOOK RECEIVED ==========")

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

	log.Println(
		"PAYPAL TRANSMISSION ID:",
		headers.TransmissionID,
	)

	body, err := io.ReadAll(
		c.Request.Body,
	)

	if err != nil {

		log.Println(
			"READ BODY FAILED:",
			err,
		)

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "cannot read body",
			},
		)

		return
	}

	log.Println(
		"BODY LENGTH:",
		len(body),
	)

	err = h.service.VerifyWebhook(
		c.Request.Context(),
		headers,
		body,
	)

	if err != nil {

		log.Println(
			"WEBHOOK VERIFICATION FAILED:",
			err,
		)

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid paypal signature",
			},
		)

		return
	}

	log.Println(
		"WEBHOOK VERIFIED",
	)

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

		log.Println(
			"JSON PARSE FAILED:",
			err,
		)

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid json",
			},
		)

		return
	}

	log.Println(
		"PAYPAL EVENT:",
		payload.EventType,
	)

	log.Println(
		"PAYPAL EVENT ID:",
		payload.EventID,
	)

	log.Println(
		"PAYPAL RESOURCE ID:",
		payload.Resource.ID,
	)

	// ---------------------------------------
	// STEP 1: Capture approved orders
	// ---------------------------------------

	if payload.EventType == "CHECKOUT.ORDER.APPROVED" {

		log.Println(
			"ORDER APPROVED - START CAPTURE",
		)

		err := h.service.CapturePayment(
			c.Request.Context(),
			payload.Resource.ID,
		)

		if err != nil {

			log.Println(
				"CAPTURE FAILED:",
				err,
			)

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "capture failed",
				},
			)

			return
		}

		log.Println(
			"CAPTURE REQUEST SUCCESSFUL",
		)

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

		log.Println(
			"IGNORING EVENT:",
			payload.EventType,
		)

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "event ignored",
			},
		)

		return
	}

	// ---------------------------------------
	// STEP 2: Idempotency
	// ---------------------------------------

	log.Println(
		"CHECKING PAYMENT EVENT:",
		payload.EventID,
	)

	existing, err := h.service.FindPaymentEvent(
		c.Request.Context(),
		payload.EventID,
	)

	if err != nil {

		log.Println(
			"PAYMENT EVENT NOT FOUND:",
			err,
		)

		existing = nil
	}

	if existing != nil {

		log.Println(
			"EXISTING EVENT FOUND. PROCESSED:",
			existing.Processed,
		)

	}

	if existing != nil && existing.Processed {

		log.Println(
			"EVENT ALREADY PROCESSED",
		)

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "event already processed",
			},
		)

		return
	}

	orderID := payload.
		Resource.
		SupplementaryData.
		RelatedIDs.
		OrderID

	log.Println(
		"CAPTURE ORDER ID:",
		orderID,
	)

	if orderID == "" {

		log.Println(
			"MISSING ORDER ID",
		)

		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "missing order id",
			},
		)

		return
	}

	// ---------------------------------------
	// STEP 3: Store event
	// ---------------------------------------

	if existing == nil {

		log.Println(
			"CREATING PAYMENT EVENT",
		)

		event := &models.PaymentEvent{

			ID: uuid.New(),

			Provider: "paypal",

			EventID: payload.EventID,

			Type: payload.EventType,

			Payload: string(body),

			Processed: false,
		}

		err = h.service.CreatePaymentEvent(
			c.Request.Context(),
			event,
		)

		if err != nil {

			log.Println(
				"CREATE PAYMENT EVENT FAILED:",
				err,
			)

			c.JSON(
				http.StatusOK,
				gin.H{
					"message": "duplicate event",
				},
			)

			return
		}

		log.Println(
			"PAYMENT EVENT STORED",
		)

	}

	// ---------------------------------------
	// STEP 4: Confirm payment
	// ---------------------------------------

	log.Println(
		"CONFIRM PAYMENT START:",
		orderID,
	)

	_, err = h.service.ConfirmPayment(
		c.Request.Context(),
		orderID,
	)

	if err != nil {

		log.Println(
			"CONFIRM PAYMENT FAILED:",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "payment processing failed",
			},
		)

		return
	}

	log.Println(
		"CONFIRM PAYMENT SUCCESS",
	)

	// ---------------------------------------
	// STEP 5: Mark webhook processed
	// ---------------------------------------

	err = h.service.MarkPaymentEventProcessed(
		c.Request.Context(),
		payload.EventID,
	)

	if err != nil {

		log.Println(
			"MARK EVENT PROCESSED FAILED:",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "could not mark event processed",
			},
		)

		return
	}

	log.Println(
		"PAYMENT EVENT MARKED PROCESSED:",
		payload.EventID,
	)

	log.Println(
		"========== PAYPAL WEBHOOK FINISHED ==========",
	)

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "payment confirmed",
		},
	)
}
