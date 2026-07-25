package scan

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestTwoStaffCannotScanSameTicketAtSameTime(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
		false,
	)

	ticketService := scenario.TicketService

	// RUN CONCURRENT SCANS

	var wg sync.WaitGroup

	wg.Add(2)

	results := make(chan error, 2)

	go func() {

		defer wg.Done()

		_, err := ticketService.Scan(
			context.Background(),
			scenario.Staff.ID,
			scenario.Ticket.Code,
		)

		results <- err
	}()

	go func() {

		defer wg.Done()

		_, err := ticketService.Scan(
			context.Background(),
			scenario.StaffTwo.ID,
			scenario.Ticket.Code,
		)

		results <- err
	}()

	wg.Wait()

	close(results)

	successes := 0

	duplicates := 0

	for err := range results {

		switch err {

		case nil:

			successes++

		case appErrors.ErrTicketAlreadyScanned:

			duplicates++

		default:

			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	}

	if successes != 1 {

		t.Fatalf(
			"expected exactly one successful scan, got %d",
			successes,
		)
	}

	if duplicates != 1 {

		t.Fatalf(
			"expected exactly one duplicate error, got %d",
			duplicates,
		)
	}

	// DATABASE CHECK

	var count int64

	if err := db.
		Model(&appModels.TicketScan{}).
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		Count(&count).
		Error; err != nil {

		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected one scan in database, got %d",
			count,
		)
	}
}

func TestSameTicketDifferentStaffSameWindow(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
		false,
	)

	ticketService := scenario.TicketService

	// FIRST STAFF SCAN

	firstScan, err := ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// SECOND STAFF SCAN

	_, err = ticketService.Scan(
		context.Background(),
		scenario.StaffTwo.ID,
		scenario.Ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected second staff scan to fail",
		)
	}

	if err != appErrors.ErrTicketAlreadyScanned {

		t.Fatalf(
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}

	// FIRST SCAN MUST BE KEPT

	if firstScan.TicketID != scenario.Ticket.ID {

		t.Fatalf(
			"expected scan for ticket %s",
			scenario.Ticket.ID,
		)
	}

	if firstScan.TicketAccessWindowID != scenario.Window.ID {

		t.Fatal(
			"expected scan in current window",
		)
	}

	// VERIFY ONLY ONE SCAN EXISTS

	var count int64

	if err := db.
		Model(&appModels.TicketScan{}).
		Where(
			"ticket_id = ? AND ticket_access_window_id = ?",
			scenario.Ticket.ID,
			scenario.Window.ID,
		).
		Count(&count).
		Error; err != nil {

		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected exactly one scan, got %d",
			count,
		)
	}
}
