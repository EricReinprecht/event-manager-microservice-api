package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/internal/requests"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPurchaseCanCreateCheckout(t *testing.T) {

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

	paypalClient := paypal.NewClient(
		"",
		"",
		"",
	)

	ticketService := helpers.NewTicketService(db)

	paymentService := helpers.NewPaymentService(
		purchaseService,
		ticketService,
		paypalClient,
	)

	url, err := paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if url == "" {

		t.Fatal(
			"expected checkout url",
		)
	}

	updated, err := purchaseService.GetPurchase(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if updated.Status != enum.StatusPending {

		t.Fatalf(
			"expected pending status got %s",
			updated.Status,
		)
	}

	if updated.PaymentProvider != "paypal" {

		t.Fatalf(
			"expected paypal provider got %s",
			updated.PaymentProvider,
		)
	}

	if updated.PaymentID == "" {

		t.Fatal(
			"expected payment id to be stored",
		)
	}
}
