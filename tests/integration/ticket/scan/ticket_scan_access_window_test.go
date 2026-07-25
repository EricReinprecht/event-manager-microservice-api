package scan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestStaffCannotScanTicketBeforeAccessWindow(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		helpers.UTCDate(
			2026,
			time.July,
			24,
			12,
		),
	)

	scenario := scenarios.CreateAccessWindowScenario(
		t,
		fakeClock,
		fakeClock.Now().Add(time.Hour),
		fakeClock.Now().Add(2*time.Hour),
	)

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err == nil {
		t.Fatal(
			"expected scan to fail before access window",
		)
	}

	if !errors.Is(err, appErrors.ErrTicketNotValidNow) {

		t.Fatalf(
			"expected ErrTicketNotValidNow, got %v",
			err,
		)
	}
}

func TestStaffCanScanTicketExactlyAtAccessWindowStart(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	scenario := scenarios.CreateAccessWindowScenario(
		t,
		fakeClock,
		fakeClock.Now().Add(time.Hour),
		fakeClock.Now().Add(2*time.Hour),
	)

	ticketService := scenario.TicketService
	ticket := scenario.Ticket
	staff := scenario.Staff

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {

		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s got %s",
			ticket.ID,
			scan.TicketID,
		)
	}
}

func TestStaffCanScanTicketExactlyAtAccessWindowEnd(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	end := fakeClock.Now().Add(time.Hour)

	scenario := scenarios.CreateAccessWindowScenario(
		t,
		fakeClock,
		end.Add(-time.Hour),
		end,
	)

	ticketService := scenario.TicketService
	ticket := scenario.Ticket
	staff := scenario.Staff

	fakeClock.Set(end)

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {

		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s got %s",
			ticket.ID,
			scan.TicketID,
		)
	}
}

func TestStaffCannotScanTicketAfterAccessWindow(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	scenario := scenarios.CreateAccessWindowScenario(
		t,
		fakeClock,
		fakeClock.Now().Add(-2*time.Hour),
		fakeClock.Now().Add(-time.Hour),
	)

	ticketService := scenario.TicketService
	ticket := scenario.Ticket
	staff := scenario.Staff

	_, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected scan to fail after access window",
		)
	}

	if !errors.Is(err, appErrors.ErrTicketNotValidNow) {

		t.Fatalf(
			"expected ErrTicketNotValidNow, got %v",
			err,
		)
	}
}

func TestScanUsesCorrectAccessWindowWhenWindowsOverlap(t *testing.T) {

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService, scenario := scenarios.CreateOverlappingAccessWindowScenario(
		t,
		clock,
	)

	scan, err := ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.TicketAccessWindowID != scenario.FirstWindow.ID {

		t.Fatalf(
			"expected scan to use first access window %s, got %s",
			scenario.FirstWindow.ID,
			scan.TicketAccessWindowID,
		)
	}
}

func TestScanTicketWithMultipleActiveWindows(t *testing.T) {

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateMultipleActiveWindowsScenario(
		t,
		clock,
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

	if scan.TicketAccessWindowID != scenario.WindowOne.ID &&
		scan.TicketAccessWindowID != scenario.WindowTwo.ID {

		t.Fatalf(
			"expected scan to belong to one active window, got %s",
			scan.TicketAccessWindowID,
		)
	}

	var count int64

	if err := scenario.DB.
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
			"expected exactly one scan, got %d",
			count,
		)
	}
}

func TestScanTicketWithoutAccessWindow(t *testing.T) {

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateNoAccessWindowScenario(
		t,
		clock,
	)

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected scan to fail without access window",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrTicketNotValidNow,
	) {

		t.Fatalf(
			"expected ErrTicketNotValidNow, got %v",
			err,
		)
	}

	var count int64

	if err := scenario.DB.
		Model(&appModels.TicketScan{}).
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		Count(&count).
		Error; err != nil {

		t.Fatal(err)
	}

	if count != 0 {

		t.Fatalf(
			"expected no ticket scans, got %d",
			count,
		)
	}
}

func TestTicketCanBeScannedInDifferentAccessWindow(t *testing.T) {

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateSequentialAccessWindowsScenario(
		t,
		clock,
	)

	firstScan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.TicketAccessWindowID != scenario.WindowOne.ID {

		t.Fatalf(
			"expected first scan in window one",
		)
	}

	// same window should reject duplicate

	_, err = scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected duplicate scan in same window",
		)
	}

	// move clock into second window

	clock.Current = clock.Current.Add(
		2 * time.Hour,
	)

	secondScan, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.TicketAccessWindowID != scenario.WindowTwo.ID {

		t.Fatalf(
			"expected second scan in window two",
		)
	}

	if secondScan.ID == firstScan.ID {

		t.Fatal(
			"expected new scan record",
		)
	}
}
