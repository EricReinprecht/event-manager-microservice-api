package integration

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

func TestNonStaffCannotScanTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

	// SCANNER USER (NOT STAFF)

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET OWNER / ORGANIZER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// NO PARTY MEMBER CREATED FOR user

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

	// ACTIVE ACCESS WINDOW

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: time.Now().
			UTC().
			Add(-time.Hour),

		EndsAt: time.Now().
			UTC().
			Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: user.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// TRY SCAN

	_, err = ticketService.Scan(
		context.Background(),
		user.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected non staff scan to fail",
		)
	}

	if err != appErrors.ErrNotAllowed {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}
}

func TestStaffFromDifferentPartyCannotScanTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

	// STAFF USER FROM ANOTHER PARTY

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY OWNER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY WITH TICKET

	partyOne := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	partyOne.CategoryID = category.ID

	if err := db.Create(&partyOne).Error; err != nil {
		t.Fatal(err)
	}

	// SECOND PARTY WHERE STAFF BELONGS

	partyTwo := fixtures.PartyWithOrganizer(
		staff.ID,
	)

	partyTwo.CategoryID = category.ID

	if err := db.Create(&partyTwo).Error; err != nil {
		t.Fatal(err)
	}

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: partyTwo.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY FROM PARTY ONE

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: partyOne.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACTIVE ACCESS WINDOW

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: time.Now().
			UTC().
			Add(-time.Hour),

		EndsAt: time.Now().
			UTC().
			Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET FROM PARTY ONE

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// STAFF FROM PARTY TWO TRIES TO SCAN PARTY ONE TICKET

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected staff from another party to be rejected",
		)
	}

	if err != appErrors.ErrNotAllowed {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}
}

func TestPartyOrganizerCanScanTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

	// ORGANIZER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER AS ORGANIZER

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: organizer.ID,

		PartyID: party.ID,

		Role: enum.RoleOrganizer,
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

		StartsAt: time.Now().
			UTC().
			Add(-time.Hour),

		EndsAt: time.Now().
			UTC().
			Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: organizer.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		organizer.ID,
		ticket.Code,
	)

	if err != nil {

		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s, got %s",
			ticket.ID,
			scan.TicketID,
		)
	}

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified scan, got %s",
			scan.Status,
		)
	}
}

func TestAdminCanScanTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

	// ADMIN USER

	admin := fixtures.User()

	if err := db.Create(&admin).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		admin.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER AS ADMIN

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: admin.ID,

		PartyID: party.ID,

		Role: enum.RoleAdmin,
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

	// ACTIVE ACCESS WINDOW

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: time.Now().
			UTC().
			Add(-time.Hour),

		EndsAt: time.Now().
			UTC().
			Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: admin.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		admin.ID,
		ticket.Code,
	)

	if err != nil {

		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s, got %s",
			ticket.ID,
			scan.TicketID,
		)
	}

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified scan, got %s",
			scan.Status,
		)
	}
}

func TestCustomerCannotScanOwnTicket(t *testing.T) {

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

	// CUSTOMER

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// STAFF / ORGANIZER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER ORGANIZER

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: organizer.ID,

		PartyID: party.ID,

		Role: enum.RoleOrganizer,
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

	// CUSTOMER TICKET

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: customer.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// CUSTOMER TRIES TO SCAN OWN TICKET

	_, err = ticketService.Scan(
		context.Background(),
		customer.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected customer scan to fail",
		)
	}

	if !errors.Is(err, appErrors.ErrNotAllowed) {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}

	// VERIFY NO SCAN WAS CREATED

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

	if count != 0 {

		t.Fatalf(
			"expected no scans, got %d",
			count,
		)
	}
}

func TestScanWithNoPartyMemberRecord(t *testing.T) {

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

	// USER (NOT A PARTY MEMBER)

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

	// REMOVE PARTY MEMBERSHIP CREATED BY FIXTURE
	// PartyWithOrganizer creates organizer membership automatically

	if err := db.
		Where(
			"user_id = ? AND party_id = ?",
			user.ID,
			party.ID,
		).
		Delete(
			&appModels.PartyMember{},
		).
		Error; err != nil {

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

	ticket.UserID = user.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// TRY SCAN WITHOUT MEMBERSHIP

	_, err = ticketService.Scan(
		context.Background(),
		user.ID,
		ticket.Code,
	)

	if err != appErrors.ErrNotAllowed {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}
}
