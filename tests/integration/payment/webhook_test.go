package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/requests"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentWebhookMarksPurchasePaid(t *testing.T) {

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

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		Price: 5000,
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

	paymentID := "PAYPAL_TEST_PAYMENT_ID"

	err = purchaseService.AttachPayment(
		context.Background(),
		purchase.ID,
		"paypal",
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	paypalClient := helpers.NewPayPalClient()

	ticketService := helpers.NewTicketService(db)

	paymentService := helpers.NewPaymentService(
		purchaseService,
		ticketService,
		paypalClient,
	)
	updated, err := paymentService.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if updated.Status != enum.StatusPaid {

		t.Fatalf(
			"expected PAID status got %s",
			updated.Status,
		)
	}

	if updated.PaymentID != paymentID {

		t.Fatalf(
			"expected payment id %s got %s",
			paymentID,
			updated.PaymentID,
		)
	}
}
