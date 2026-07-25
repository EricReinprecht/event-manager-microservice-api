package scan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestScanTransactionRollback(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		staff.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: time.Now().UTC().Add(-time.Hour),

		EndsAt: time.Now().UTC().Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = staff.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// first scan succeeds

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// second scan must fail

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected duplicate scan error",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrTicketAlreadyScanned,
	) {

		t.Fatalf(
			"expected duplicate error, got %v",
			err,
		)
	}

	// ensure only one record exists

	var count int64

	if err := db.
		Model(&appModels.TicketScan{}).
		Where(
			"ticket_id = ?",
			ticket.ID,
		).
		Count(&count).
		Error; err != nil {

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

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().UTC(),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// STAFF

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

	// MEMBER

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

	// ACCESS WINDOW

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = staff.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// CREATE EXISTING PENDING SCAN IN SAME WINDOW

	existingScan := appModels.TicketScan{

		ID: uuid.New(),

		TicketID: ticket.ID,

		TicketAccessWindowID: window.ID,

		ScannedByID: staff.ID,

		ScannedAt: clock.Current,

		Status: enum.TicketScanPending,
	}

	if err := db.Create(&existingScan).Error; err != nil {
		t.Fatal(err)
	}

	// TRY SECOND SCAN

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != appErrors.ErrTicketAlreadyScanned {

		t.Fatalf(
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}

	// VERIFY NO ADDITIONAL SCAN WAS CREATED

	var count int64

	if err := db.Model(
		&appModels.TicketScan{},
	).
		Where(
			"ticket_id = ?",
			ticket.ID,
		).
		Count(&count).
		Error; err != nil {

		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected exactly one scan after rollback, got %d",
			count,
		)
	}
}
