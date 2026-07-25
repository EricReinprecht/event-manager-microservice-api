package paypal

import (
	"context"
	"testing"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentCapturePaymentFailed(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

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

	err = paymentService.CapturePayment(
		context.Background(),
		"FAILED_ORDER_ID",
	)

	if err == nil {

		t.Fatal(
			"expected capture error",
		)
	}
}
