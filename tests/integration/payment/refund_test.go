package payment

import (
	"context"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPaymentRefundSuccess(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
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
		enum.PurchaseStatusPaid,
	)

	purchase.PaymentProvider = "paypal"

	purchase.PaymentID = "PAYPAL-PAYMENT-123"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	err = paymentService.RefundPayment(
		context.Background(),
		purchase.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

	if fakeGateway.RefundedPaymentID != purchase.PaymentID {

		t.Fatalf(
			"expected refund payment id %s got %s",
			purchase.PaymentID,
			fakeGateway.RefundedPaymentID,
		)
	}

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusRefunded {

		t.Fatalf(
			"expected REFUNDED got %s",
			updated.Status,
		)
	}

	if updated.RefundID == "" {

		t.Fatal(
			"expected refund id to be stored",
		)
	}

	if updated.RefundedAt == nil {

		t.Fatal(
			"expected refunded timestamp",
		)
	}
}

func TestPaymentRefundGatewayFailure(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	failingGateway := &helpers.FailingPaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		failingGateway,
	)

	// ----------------------------
	// Setup user
	// ----------------------------

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

	// ----------------------------
	// Create paid purchase
	// ----------------------------

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPaid,
	)

	purchase.PaymentProvider = "paypal"

	purchase.PaymentID = "PAYPAL-FAIL-REFUND"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// ----------------------------
	// Execute refund
	// ----------------------------

	err = paymentService.RefundPayment(
		context.Background(),
		purchase.ID,
	)

	if err == nil {

		t.Fatal(
			"expected refund failure",
		)
	}

	// ----------------------------
	// Verify purchase unchanged
	// ----------------------------

	var updated models.Purchase

	if err := db.
		First(
			&updated,
			"id = ?",
			purchase.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase to remain PAID, got %s",
			updated.Status,
		)
	}

	if updated.RefundID != "" {

		t.Fatalf(
			"expected no refund id, got %s",
			updated.RefundID,
		)
	}

	if updated.RefundedAt != nil {

		t.Fatal(
			"expected refunded timestamp to remain nil",
		)
	}
}

func TestCannotRefundAlreadyRefundedPurchase(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
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

	// USER

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// ALREADY REFUNDED PURCHASE

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.PurchaseStatusPaid,
	)

	purchase.PaymentProvider = "paypal"

	purchase.PaymentID = "PAYPAL-ALREADY-REFUNDED"

	purchase.Status = enum.PurchaseStatusRefunded

	purchase.RefundID = "REFUND-123"

	purchase.RefundProvider = "paypal"

	now := time.Now().UTC()

	purchase.RefundedAt = &now

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// TRY REFUND AGAIN

	err = paymentService.RefundPayment(
		context.Background(),
		purchase.ID,
	)

	if err == nil {

		t.Fatal(
			"expected already refunded purchase to fail",
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"refund gateway should not be called twice",
		)
	}
}
