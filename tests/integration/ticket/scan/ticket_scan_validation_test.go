package scan

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestCannotScanUnknownTicketCode(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := helpers.NewFakeClock(time.Now().UTC())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
		false,
	)

	ticketService := scenario.TicketService

	_, err = ticketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		"THIS-TICKET-DOES-NOT-EXIST",
	)

	if err == nil {

		t.Fatal(
			"expected unknown ticket code error",
		)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) &&
		!errors.Is(err, appErrors.ErrTicketNotFound) {

		t.Fatalf(
			"expected ticket not found error, got %v",
			err,
		)
	}
}

func TestCannotScanDeletedTicket(t *testing.T) {

	db, _ := helpers.TestDatabase()

	helpers.CleanDatabase(db)

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		&helpers.FakeClock{
			Current: time.Now().UTC(),
		},
	)

	code := scenario.Ticket.Code

	if err := db.Delete(
		&scenario.Ticket,
	).Error; err != nil {
		t.Fatal(err)
	}

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		code,
	)

	if !errors.Is(err, gorm.ErrRecordNotFound) {

		t.Fatalf(
			"expected not found, got %v",
			err,
		)
	}
}

func TestCannotScanTicketWithDeletedCategory(t *testing.T) {

	db, _ := helpers.TestDatabase()

	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(time.Now().UTC())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	if err := db.Delete(
		&scenario.TicketCategory,
	).Error; err != nil {
		t.Fatal(err)
	}

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(err, gorm.ErrRecordNotFound) {

		t.Fatalf(
			"expected category missing, got %v",
			err,
		)
	}
}

func TestTicketCannotBeScannedWithoutActivePartyMembership(t *testing.T) {

	db, _ := helpers.TestDatabase()

	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(time.Now().UTC())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	scenarios.RemovePartyMember(
		t,
		db,
		&scenario.Member,
	)

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed got %v",
			err,
		)
	}
}

func TestStaffCannotScanCancelledTicket(t *testing.T) {

	db, _ := helpers.TestDatabase()

	helpers.CleanDatabase(db)

	clock := helpers.NewFakeClock(time.Now().UTC())

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
	)

	scenarios.CancelTicket(
		t,
		db,
		&scenario.Ticket,
	)

	_, err := scenario.TicketService.Scan(
		context.Background(),
		scenario.Staff.ID,
		scenario.Ticket.Code,
	)

	if !errors.Is(
		err,
		appErrors.ErrTicketNotValidNow,
	) {

		t.Fatalf(
			"expected invalid ticket error got %v",
			err,
		)
	}

	var count int64

	db.Model(
		&models.TicketScan{},
	).
		Where(
			"ticket_id = ?",
			scenario.Ticket.ID,
		).
		Count(&count)

	if count != 0 {

		t.Fatalf(
			"expected no scans, got %d",
			count,
		)
	}
}
