package scan

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestDatabasePreventsDuplicateTicketWindowScan(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := helpers.NewFakeClock(
		time.Now().
			UTC().
			Truncate(time.Microsecond),
	)

	scenario := scenarios.CreateScanScenario(
		t,
		db,
		clock,
		false,
	)

	now := clock.Now()

	firstScan := appModels.TicketScan{

		ID: uuid.New(),

		TicketID: scenario.Ticket.ID,

		TicketAccessWindowID: scenario.Window.ID,

		ScannedByID: scenario.Staff.ID,

		Status: enum.TicketScanPending,

		ScannedAt: now,
	}

	if err := db.Create(&firstScan).Error; err != nil {
		t.Fatal(err)
	}

	secondScan := firstScan

	secondScan.ID = uuid.New()

	secondScan.ScannedAt = now.Add(time.Second)

	err = db.Create(&secondScan).Error

	if err == nil {

		t.Fatal(
			"expected database duplicate constraint error",
		)
	}

	if !strings.Contains(
		err.Error(),
		"idx_ticket_window_pending",
	) {

		t.Fatalf(
			"expected idx_ticket_window_pending constraint error, got %v",
			err,
		)
	}
}
