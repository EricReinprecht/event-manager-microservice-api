package scan

import (
	"context"
	"testing"
	"time"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestStaffCanScanTicketWithoutVerification(t *testing.T) {

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

	result, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.TicketID != scenario.Ticket.ID {

		t.Fatalf(
			"wrong ticket scanned expected %s got %s",
			scenario.Ticket.ID,
			result.TicketID,
		)
	}

	var scan appModels.TicketScan

	err = db.
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		First(&scan).
		Error

	if err != nil {

		t.Fatal(
			"ticket scan was not created",
		)
	}

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified got %s",
			scan.Status,
		)
	}

	if scan.VerifiedAt == nil {

		t.Fatal(
			"expected verified_at",
		)
	}

	if scan.VerifiedByID == nil {

		t.Fatal(
			"expected verified_by_id",
		)
	}

	if *scan.VerifiedByID != scenario.Staff.ID {

		t.Fatalf(
			"expected verifier %s got %s",
			scenario.Staff.ID,
			*scan.VerifiedByID,
		)
	}
}

func TestTwoDifferentTicketsCanBeScannedBySameStaff(t *testing.T) {

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

	customerTwo := fixtures.User()

	if err := db.Create(&customerTwo).Error; err != nil {
		t.Fatal(err)
	}

	ticketTwo := scenarios.CreateAdditionalTicket(
		t,
		db,
		scenario.TicketCategory.ID,
		customerTwo.ID,
		scenario.Party.ID,
	)

	firstScan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	secondScan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		ticketTwo.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.ScannedByID != scenario.Staff.ID {

		t.Fatal(
			"first scan wrong staff",
		)
	}

	if secondScan.ScannedByID != scenario.Staff.ID {

		t.Fatal(
			"second scan wrong staff",
		)
	}

	if firstScan.TicketAccessWindowID != scenario.Window.ID {

		t.Fatal(
			"first scan wrong window",
		)
	}

	if secondScan.TicketAccessWindowID != scenario.Window.ID {

		t.Fatal(
			"second scan wrong window",
		)
	}
}
