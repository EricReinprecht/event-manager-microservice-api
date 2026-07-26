package paypal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentCreateCheckoutPayPalUnavailable(t *testing.T) {

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

	fakeGateway := &helpers.FailingPaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

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

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPending,
	)

	_, err = paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
	)

	if err == nil {

		t.Fatal(
			"expected PayPal unavailable error",
		)
	}

	var updated models.Purchase

	err = db.
		First(
			&updated,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updated.PaymentID != "" {

		t.Fatal(
			"expected payment id not to be saved when PayPal fails",
		)
	}
}

func TestPayPalCreateOrderSuccess(t *testing.T) {

	err := helpers.LoadTestEnv()

	if err != nil {
		t.Fatal(err)
	}

	client := helpers.NewPayPalClient()

	order, err := client.CreateOrder(
		context.Background(),
		1000,
	)

	if err != nil {
		t.Fatalf(
			"paypal order creation failed: %v",
			err,
		)
	}

	if order.ID == "" {
		t.Fatal("expected paypal order id")
	}

	if order.ApprovalURL == "" {
		t.Fatal("expected approval url")
	}
}

func TestPayPalCreateOrderInvalidAmount(t *testing.T) {

	err := helpers.LoadTestEnv()

	if err != nil {
		t.Fatal(err)
	}

	client := helpers.NewPayPalClient()

	_, err = client.CreateOrder(
		context.Background(),
		0,
	)

	if err == nil {
		t.Fatal(
			"expected error for invalid amount",
		)
	}
}

func TestPayPalCreateOrderAPIError(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				w.WriteHeader(
					http.StatusInternalServerError,
				)

				w.Write([]byte(`{
					"error": "server_error"
				}`))
			},
		),
	)

	defer server.Close()

	client := paypal.NewClientWithBaseURL(
		"test",
		"test",
		server.URL,
		"http://return",
		"http://cancel",
		"webhook",
	)

	_, err := client.CreateOrder(
		context.Background(),
		1000,
	)

	if err == nil {
		t.Fatal(
			"expected paypal api error",
		)
	}
}

func TestPayPalCaptureOrderSuccess(t *testing.T) {

	err := helpers.LoadTestEnv()

	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)

	defer cancel()

	client := helpers.NewPayPalClient()

	// Create PayPal order

	order, err := client.CreateOrder(
		ctx,
		10,
	)

	if err != nil {

		t.Fatalf(
			"create order failed: %v",
			err,
		)
	}

	if order.ID == "" {

		t.Fatal(
			"paypal order id is empty",
		)
	}

	if order.ApprovalURL == "" {

		t.Fatal(
			"paypal approval url is empty",
		)
	}

	// Approve order as sandbox buyer

	err = helpers.ApprovePayPalOrder(
		order.ApprovalURL,
	)

	if err != nil {

		t.Fatalf(
			"paypal approval failed: %v",
			err,
		)
	}

	// Give PayPal sandbox time to update order state

	time.Sleep(
		2 * time.Second,
	)

	// Capture payment

	err = client.CaptureOrder(
		ctx,
		order.ID,
	)

	if err != nil {

		t.Fatalf(
			"capture failed for order %s: %v",
			order.ID,
			err,
		)
	}
}
