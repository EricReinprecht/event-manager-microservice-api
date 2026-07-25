package verification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestVerificationFailsExactlyAtExpiryBoundary(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scenario := scenarios.CreateVerificationScenario(
		t,
		db,
		clock,
		true,
	)

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	scenarios.AssertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	if scan.VerificationExpiresAt == nil {
		t.Fatal(
			"expected verification expiry to be set",
		)
	}

	// MOVE EXACTLY TO EXPIRY TIME

	clock.Current = *scan.VerificationExpiresAt

	// VERIFY AT BOUNDARY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		scenario.Staff.ID,
		true,
	)

	if err != nil {
		t.Fatalf(
			"expected verification to succeed at expiry boundary, got %v",
			err,
		)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	scenarios.AssertScanStatus(
		t,
		&updated,
		enum.TicketScanVerified,
	)

	if updated.VerifiedAt == nil {
		t.Fatal(
			"expected VerifiedAt to be set",
		)
	}

	if updated.VerifiedByID == nil {
		t.Fatal(
			"expected VerifiedByID to be set",
		)
	}
}

func TestPendingScanExpiresCannotBeVerified(t *testing.T) {

	db, err := helpers.TestDatabase()
	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scenario := scenarios.CreateVerificationScenario(
		t,
		db,
		clock,
		true,
	)

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	scenarios.AssertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	if scan.VerificationExpiresAt == nil {
		t.Fatal(
			"expected verification expiry to be set",
		)
	}

	// MOVE AFTER EXPIRY

	clock.Current = scan.VerificationExpiresAt.Add(
		time.Second,
	)

	// VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		scenario.Staff.ID,
		true,
	)

	if !errors.Is(
		err,
		appErrors.ErrTicketVerificationExpired,
	) {

		t.Fatalf(
			"expected verification expired error, got %v",
			err,
		)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	scenarios.AssertScanStatus(
		t,
		&updated,
		enum.TicketScanPending,
	)

	if updated.VerifiedAt != nil {

		t.Fatal(
			"expected VerifiedAt to remain nil",
		)
	}

	if updated.VerifiedByID != nil {

		t.Fatal(
			"expected VerifiedByID to remain nil",
		)
	}
}

func TestVerificationExpiredPendingScanCannotBeVerified(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scenario := scenarios.CreateVerificationScenario(
		t,
		db,
		clock,
		true,
	)

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	scenarios.AssertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	if scan.VerificationExpiresAt == nil {
		t.Fatal(
			"expected verification expiry to be set",
		)
	}

	// MOVE BEYOND VERIFICATION TTL

	clock.Current = clock.Current.Add(
		30 * time.Minute,
	)

	// TRY VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		scenario.Staff.ID,
		true,
	)

	if err == nil {
		t.Fatal(
			"expected expired verification error",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrTicketVerificationExpired,
	) {

		t.Fatalf(
			"expected ErrTicketVerificationExpired, got %v",
			err,
		)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	scenarios.AssertScanStatus(
		t,
		&updated,
		enum.TicketScanPending,
	)
}
