package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestTicketCanBeScannedInDifferentAccessWindow(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	fakeClock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	ticketService := helpers.NewTicketService(
		db,
		fakeClock,
	)

	// USER

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		staff.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER

	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&member).Error; err != nil {
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

	// FIRST ACCESS WINDOW

	windowOne := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: fakeClock.Current.Add(-time.Hour),

		EndsAt: fakeClock.Current.Add(time.Hour),
	}

	if err := db.Create(&windowOne).Error; err != nil {
		t.Fatal(err)
	}

	// SECOND ACCESS WINDOW

	windowTwo := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: fakeClock.Current.Add(2 * time.Hour),

		EndsAt: fakeClock.Current.Add(3 * time.Hour),
	}

	if err := db.Create(&windowTwo).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// FIRST SCAN IN WINDOW ONE

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.TicketAccessWindowID != windowOne.ID {

		t.Fatalf(
			"expected first scan in window one",
		)
	}

	// SECOND SCAN IN SAME WINDOW SHOULD FAIL

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected duplicate scan in same window",
		)
	}

	// MOVE CLOCK INTO SECOND WINDOW

	fakeClock.Current = fakeClock.Current.Add(
		2 * time.Hour,
	)

	// SCAN AGAIN IN SECOND WINDOW

	secondScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.TicketAccessWindowID != windowTwo.ID {

		t.Fatalf(
			"expected second scan in window two",
		)
	}

	if secondScan.ID == firstScan.ID {

		t.Fatal(
			"expected a new scan record",
		)
	}
}
