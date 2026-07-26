package refund

import (
	"context"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestPaymentRefundSuccess(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.SetupPaymentTestService(
		db,
		fakeGateway,
	)
	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	// -------------------------
	// Execute refund
	// -------------------------

	err = paymentService.RefundPayment(
		context.Background(),
		scenario.Purchase.ID,
		scenario.Actor.ID,
	)

	if err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Gateway assertions
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
	// Database assertions
	// -------------------------

	var updated models.Purchase

	if err := db.
		First(
			&updated,
			"id = ?",
			scenario.Purchase.ID,
		).
		Error; err != nil {

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

	failingGateway := &helpers.FailingPaymentGateway{}

	paymentService := helpers.SetupPaymentTestService(
		db,
		failingGateway,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	err = paymentService.RefundPayment(
		context.Background(),
		scenario.Purchase.ID,
		scenario.Actor.ID,
	)

	if err == nil {
		t.Fatal(
			"expected refund failure",
		)
	}

	if err.Error() != "refund failed" {

		t.Fatalf(
			"expected refund failed error, got %v",
			err,
		)
	}

	var updated models.Purchase

	if err := db.
		First(
			&updated,
			"id = ?",
			scenario.Purchase.ID,
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

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.SetupPaymentTestService(
		db,
		fakeGateway,
	)

	// ----------------------------
	// Scenario
	// ----------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	purchase := scenario.Purchase

	purchase.Status = enum.PurchaseStatusRefunded
	purchase.RefundID = "REFUND-123"
	purchase.RefundProvider = "paypal"

	now := time.Now().UTC()

	purchase.RefundedAt = &now

	if err := db.Save(
		purchase,
	).Error; err != nil {

		t.Fatal(err)
	}

	// ----------------------------
	// Execute refund again
	// ----------------------------

	err = paymentService.RefundPayment(
		context.Background(),
		purchase.ID,
		scenario.Actor.ID,
	)

	if err == nil {

		t.Fatal(
			"expected already refunded purchase to fail",
		)
	}

	// ----------------------------
	// Gateway must not be touched
	// ----------------------------

	if fakeGateway.RefundCalled {

		t.Fatal(
			"refund gateway should not be called twice",
		)
	}

	// ----------------------------
	// Verify unchanged
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

	if updated.Status != enum.PurchaseStatusRefunded {

		t.Fatalf(
			"expected REFUNDED status, got %s",
			updated.Status,
		)
	}

	if updated.RefundID != "REFUND-123" {

		t.Fatalf(
			"expected existing refund id to remain, got %s",
			updated.RefundID,
		)
	}
}
