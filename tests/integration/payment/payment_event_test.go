package payment

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentEvent_CreatesEvent(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	eventID := "WH-CREATE-EVENT-" + uuid.New().String()

	body := fmt.Appendf(
		nil,
		`
{
	"id": "%s",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "ORDER-123"

			}
		}
	}
}
`,
		eventID,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"test-id",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"test-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	// event creation happens before payment confirmation,
	// therefore the payment will fail because ORDER-123
	// does not exist. We only care about the event.

	if w.Code != http.StatusInternalServerError {

		t.Fatalf(
			"expected 500 because order is missing, got %d body %s",
			w.Code,
			w.Body.String(),
		)
	}

	var event models.PaymentEvent

	err = db.
		Where(
			"event_id = ?",
			eventID,
		).
		First(&event).
		Error

	if err != nil {

		t.Fatal(
			"expected payment event to be created:",
			err,
		)
	}

	if event.Provider != "paypal" {

		t.Fatalf(
			"expected provider paypal, got %s",
			event.Provider,
		)
	}

	if event.EventID != eventID {

		t.Fatalf(
			"expected event id %s, got %s",
			eventID,
			event.EventID,
		)
	}

	if event.Type != "PAYMENT.CAPTURE.COMPLETED" {

		t.Fatalf(
			"expected event type PAYMENT.CAPTURE.COMPLETED, got %s",
			event.Type,
		)
	}

	if event.Processed {

		t.Fatal(
			"expected event to be unprocessed",
		)
	}
}

func TestPaymentEvent_EventIDMustBeUnique(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	eventID := "WH-UNIQUE-" + uuid.New().String()

	firstEvent := models.PaymentEvent{

		ID: uuid.New(),

		Provider: "paypal",

		EventID: eventID,

		Type: "PAYMENT.CAPTURE.COMPLETED",

		Payload: "{}",

		Processed: false,
	}

	err = db.Create(&firstEvent).Error

	if err != nil {
		t.Fatal(
			"first event should be created:",
			err,
		)
	}

	secondEvent := models.PaymentEvent{

		ID: uuid.New(),

		Provider: "paypal",

		EventID: eventID,

		Type: "PAYMENT.CAPTURE.COMPLETED",

		Payload: "{}",

		Processed: false,
	}

	err = db.Create(&secondEvent).Error

	if err == nil {

		t.Fatal(
			"expected duplicate event id error",
		)
	}
}
