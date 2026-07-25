package scan

import (
	"context"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestScanCreatesCorrectScannerMetadata(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(
		helpers.UTCDate(
			2026,
			time.July,
			24,
			12,
		),
	)

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
		false,
	)

	scan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.TicketID != scenario.Ticket.ID {
		t.Fatalf(
			"expected ticket %s got %s",
			scenario.Ticket.ID,
			scan.TicketID,
		)
	}

	if scan.ScannedByID != scenario.Staff.ID {
		t.Fatalf(
			"expected scanner %s got %s",
			scenario.Staff.ID,
			scan.ScannedByID,
		)
	}

	if !scan.ScannedAt.Equal(clock.Now()) {
		t.Fatalf(
			"expected scanned at %v got %v",
			clock.Now(),
			scan.ScannedAt,
		)
	}

	if scan.TicketAccessWindowID != scenario.Window.ID {
		t.Fatalf(
			"expected window %s got %s",
			scenario.Window.ID,
			scan.TicketAccessWindowID,
		)
	}

	if scan.Status != enum.TicketScanVerified {
		t.Fatalf(
			"expected verified got %s",
			scan.Status,
		)
	}

	if scan.VerifiedByID == nil {
		t.Fatal(
			"expected verifier",
		)
	}

	if *scan.VerifiedByID != scenario.Staff.ID {
		t.Fatalf(
			"expected verifier %s got %s",
			scenario.Staff.ID,
			*scan.VerifiedByID,
		)
	}

	if scan.VerifiedAt == nil {
		t.Fatal(
			"expected verified timestamp",
		)
	}

	if !scan.VerifiedAt.Equal(clock.Now()) {
		t.Fatalf(
			"expected verified at %v got %v",
			clock.Now(),
			*scan.VerifiedAt,
		)
	}
}

func TestVerifiedScanStoresVerificationMetadata(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(
		time.Now().UTC(),
	)

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
		false,
	)

	scan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanVerified {
		t.Fatalf(
			"expected verified got %s",
			scan.Status,
		)
	}

	if scan.VerifiedAt == nil {
		t.Fatal(
			"expected VerifiedAt",
		)
	}

	if !scan.VerifiedAt.Equal(clock.Now()) {
		t.Fatalf(
			"expected %v got %v",
			clock.Now(),
			*scan.VerifiedAt,
		)
	}

	if scan.VerifiedByID == nil {
		t.Fatal(
			"expected VerifiedByID",
		)
	}

	if *scan.VerifiedByID != scenario.Staff.ID {
		t.Fatalf(
			"expected %s got %s",
			scenario.Staff.ID,
			*scan.VerifiedByID,
		)
	}

	expectedExpiry := clock.Now().Add(
		15 * time.Minute,
	)

	if scan.VerifiedUntil == nil {
		t.Fatal(
			"expected VerifiedUntil",
		)
	}

	if !scan.VerifiedUntil.Equal(expectedExpiry) {
		t.Fatalf(
			"expected %v got %v",
			expectedExpiry,
			*scan.VerifiedUntil,
		)
	}
}
