package refund

import (
	"context"
	"testing"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestOrganizerCanRefundPurchase(t *testing.T) {

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

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	scenarios.RefundPurchase(
		t,
		paymentService,
		scenario.Purchase.ID,
		scenario.Organizer.ID,
	)

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

	if fakeGateway.RefundedPaymentID != scenario.Purchase.PaymentID {

		t.Fatalf(
			"expected refund payment id %s got %s",
			scenario.Purchase.PaymentID,
			fakeGateway.RefundedPaymentID,
		)
	}

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
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
			"expected refund id",
		)
	}

	if updated.RefundProvider != "paypal" {

		t.Fatalf(
			"expected paypal refund provider got %s",
			updated.RefundProvider,
		)
	}

	if updated.RefundedAt == nil {

		t.Fatal(
			"expected refunded timestamp",
		)
	}
}

func TestAdminCanRefundPurchase(t *testing.T) {

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

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAdmin,
	)

	// -------------------------
	// Execute refund as admin
	// -------------------------

	scenarios.RefundPurchase(
		t,
		paymentService,
		scenario.Purchase.ID,
		scenario.Actor.ID,
	)

	// -------------------------
	// Verify gateway
	// -------------------------

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

	if fakeGateway.RefundedPaymentID != scenario.Purchase.PaymentID {

		t.Fatalf(
			"expected refund payment id %s got %s",
			scenario.Purchase.PaymentID,
			fakeGateway.RefundedPaymentID,
		)
	}

	// -------------------------
	// Verify purchase updated
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
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
			"expected refund id",
		)
	}

	if updated.RefundProvider != "paypal" {

		t.Fatalf(
			"expected paypal refund provider got %s",
			updated.RefundProvider,
		)
	}

	if updated.RefundedAt == nil {

		t.Fatal(
			"expected refunded timestamp",
		)
	}
}

func TestRefunderRoleCanRefundPurchase(t *testing.T) {

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

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleRefunder,
	)

	// -------------------------
	// Execute refund
	// -------------------------

	scenarios.RefundPurchase(
		t,
		paymentService,
		scenario.Purchase.ID,
		scenario.Actor.ID,
	)

	// -------------------------
	// Verify gateway
	// -------------------------

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

	if fakeGateway.RefundedPaymentID != scenario.Purchase.PaymentID {

		t.Fatalf(
			"expected refund payment id %s got %s",
			scenario.Purchase.PaymentID,
			fakeGateway.RefundedPaymentID,
		)
	}

	// -------------------------
	// Verify database
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
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
			"expected refund id",
		)
	}

	if updated.RefundProvider != "paypal" {

		t.Fatalf(
			"expected paypal refund provider got %s",
			updated.RefundProvider,
		)
	}

	if updated.RefundedAt == nil {

		t.Fatal(
			"expected refunded timestamp",
		)
	}
}

func TestStaffCannotRefundPurchase(t *testing.T) {

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

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleStaff,
	)

	err = paymentService.RefundPayment(
		context.Background(),
		scenario.Purchase.ID,
		scenario.Actor.ID,
	)

	if err == nil {

		t.Fatal(
			"expected staff refund to be rejected",
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"staff should not reach refund gateway",
		)
	}

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase to remain PAID, got %s",
			updated.Status,
		)
	}
}

func TestAttendeeCannotRefundPurchase(t *testing.T) {

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

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAttendee,
	)

	// -------------------------
	// Try refund as attendee
	// -------------------------

	err = paymentService.RefundPayment(
		context.Background(),
		scenario.Purchase.ID,
		scenario.Actor.ID,
	)

	if err == nil {

		t.Fatal(
			"expected attendee refund to be rejected",
		)
	}

	// -------------------------
	// Verify gateway was blocked
	// -------------------------

	if fakeGateway.RefundCalled {

		t.Fatal(
			"attendee should not reach refund gateway",
		)
	}

	// -------------------------
	// Verify purchase unchanged
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected purchase to remain PAID, got %s",
			updated.Status,
		)
	}
}
