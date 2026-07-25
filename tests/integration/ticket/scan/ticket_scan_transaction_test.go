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

func TestScanTransactionRollback(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// first scan succeeds

	_, err = ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// second scan fails

	_, err = ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrTicketAlreadyScanned,
	) {

		t.Fatalf(
			"expected duplicate error, got %v",
			err,
		)
	}

	var count int64

	err = db.
		Model(&appModels.TicketScan{}).
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		Count(&count).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected 1 scan after rollback, got %d",
			count,
		)
	}
}

func TestScanRollbackAfterDatabaseConstraintFailure(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	scenarios.CreatePendingTicketScan(
		t,
		db,
		scenario.Ticket,
		scenario.Window,
		scenario.Staff,
		clock.Current,
	)

	_, err = ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrTicketAlreadyScanned,
	) {

		t.Fatalf(
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}

	var count int64

	err = db.
		Model(&appModels.TicketScan{}).
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		Count(&count).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected exactly one scan after rollback, got %d",
			count,
		)
	}
}
