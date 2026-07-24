package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/handlers"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	paypalClient "github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/requests"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentWebhookConfirmsPurchase(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := appModels.TicketCategory{
		ID:      uuid.New(),
		Name:    "VIP",
		PartyID: party.ID,
		Price:   5000,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	purchaseService := helpers.NewPurchaseService(db)

	purchase, err := purchaseService.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: ticketCategory.ID,
				Quantity:         1,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	paymentID := "PAYPAL_TEST_ORDER_ID"

	err = purchaseService.AttachPayment(
		context.Background(),
		purchase.ID,
		"paypal",
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	paypal := paypalClient.NewClient(
		"",
		"",
		"",
	)

	ticketService := helpers.NewTicketService(db)

	paymentService := helpers.NewPaymentService(
		purchaseService,
		ticketService,
		paypal,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.POST(
		"/api/payments/paypal/webhook",
		handler.Webhook,
	)

	payload := map[string]interface{}{
		"id": "WH-TEST",
		"resource": map[string]interface{}{
			"id": paymentID,
		},
	}

	body, err := json.Marshal(payload)

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/payments/paypal/webhook",
		bytes.NewBuffer(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"test",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"test",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"test",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"test",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {

		t.Fatalf(
			"expected status 200 got %d body %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	updated, err := purchaseService.GetPurchase(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if updated.Status != enum.StatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updated.Status,
		)
	}

	tickets := []appModels.Ticket{}

	if err := db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).Error; err != nil {

		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 ticket got %d",
			len(tickets),
		)
	}
}
