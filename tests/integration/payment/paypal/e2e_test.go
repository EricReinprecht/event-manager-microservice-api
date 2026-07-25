package paypal

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestE2E_PayPalSandbox_PaymentFlow(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	paypalClient := helpers.NewPayPalClient()

	paymentService := helpers.NewPaymentService(
		database.NewGormExecutor(db),
		purchaseService,
		ticketService,
		paypalClient,
	)

	// --------------------------------
	// Create test purchase
	// --------------------------------

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PruchaseStatusPending,
	)

	// --------------------------------
	// Create ticket data
	// --------------------------------

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	ticketCategory.Price = 1000

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	// --------------------------------
	// Create PayPal checkout
	// --------------------------------

	var items []models.PurchaseItem

	if err := db.
		Where(
			"purchase_id = ?",
			purchase.ID,
		).
		Find(&items).
		Error; err != nil {
		t.Fatal(err)
	}

	for _, item := range items {
		t.Logf(
			"purchase item: category=%s quantity=%d price=%d",
			item.TicketCategoryID,
			item.Quantity,
			item.UnitPrice,
		)
	}

	approveURL, err := paymentService.CreateCheckout(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if approveURL == "" {

		t.Fatal(
			"paypal approve url missing",
		)
	}

	t.Log(
		"Approve URL:",
		approveURL,
	)

	// --------------------------------
	// Approve order with browser
	// --------------------------------

	err = helpers.ApprovePayPalOrder(
		approveURL,
	)

	if err != nil {
		t.Fatalf(
			"paypal approval failed: %v",
			err,
		)
	}

	// --------------------------------
	// Reload purchase to get PayPal order ID
	// --------------------------------

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.PaymentID == "" {

		t.Fatal(
			"paypal order id was not saved",
		)
	}

	// --------------------------------
	// Capture payment
	// --------------------------------

	err = paymentService.CapturePayment(
		context.Background(),
		updatedPurchase.PaymentID,
	)

	if err != nil {
		t.Fatalf(
			"paypal capture failed for order %s: %v",
			updatedPurchase.PaymentID,
			err,
		)
	}

	// --------------------------------
	// Simulate completed webhook
	// --------------------------------

	_, err = paymentService.ConfirmPayment(
		context.Background(),
		updatedPurchase.PaymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	// --------------------------------
	// Verify purchase paid
	// --------------------------------

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updatedPurchase.Status,
		)
	}

	// --------------------------------
	// Verify tickets generated
	// --------------------------------

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 ticket, got %d",
			len(tickets),
		)
	}
}
