package scan

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestDatabasePreventsDuplicateTicketWindowScan(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	// USER

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.UserID = user.ID

	ticket.TicketCategoryID = ticketCategory.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	window := fixtures.TicketAccessWindow(
		ticketCategory.ID,
	)

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// FIRST ACTIVE SCAN

	firstScan := appModels.TicketScan{

		ID: uuid.New(),

		TicketID: ticket.ID,

		TicketAccessWindowID: window.ID,

		ScannedByID: user.ID,

		Status: enum.TicketScanPending,

		ScannedAt: time.Now().UTC(),
	}

	if err := db.Create(&firstScan).Error; err != nil {
		t.Fatal(err)
	}

	// SECOND ACTIVE SCAN

	secondScan := appModels.TicketScan{

		ID: uuid.New(),

		TicketID: ticket.ID,

		TicketAccessWindowID: window.ID,

		ScannedByID: user.ID,

		Status: enum.TicketScanPending,

		ScannedAt: time.Now().UTC(),
	}

	err = db.Create(&secondScan).Error

	// EXPECT DATABASE CONSTRAINT FAILURE

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
