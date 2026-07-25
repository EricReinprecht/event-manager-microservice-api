package paypal

import (
	"context"
	"testing"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
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
		enum.StatusPending,
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
