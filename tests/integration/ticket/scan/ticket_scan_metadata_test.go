package scan

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

func TestScanCreatesCorrectScannerMetadata(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
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

	// CUSTOMER

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
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

	// PARTY STAFF

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

	// WINDOW

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
	ticket.UserID = customer.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// ASSERT METADATA

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket id %s, got %s",
			ticket.ID,
			scan.TicketID,
		)
	}

	if scan.ScannedByID != staff.ID {

		t.Fatalf(
			"expected scanner id %s, got %s",
			staff.ID,
			scan.ScannedByID,
		)
	}

	if !scan.ScannedAt.Equal(clock.Current) {

		t.Fatalf(
			"expected scanned at %v, got %v",
			clock.Current,
			scan.ScannedAt,
		)
	}

	if scan.TicketAccessWindowID != window.ID {

		t.Fatalf(
			"expected access window %s, got %s",
			window.ID,
			scan.TicketAccessWindowID,
		)
	}

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified status, got %s",
			scan.Status,
		)
	}

	if scan.VerifiedByID == nil {

		t.Fatal(
			"expected verified by id",
		)
	}

	if *scan.VerifiedByID != staff.ID {

		t.Fatalf(
			"expected verifier %s, got %s",
			staff.ID,
			*scan.VerifiedByID,
		)
	}

	if scan.VerifiedAt == nil {

		t.Fatal(
			"expected verified timestamp",
		)
	}

	if !scan.VerifiedAt.Equal(clock.Current) {

		t.Fatalf(
			"expected verified at %v, got %v",
			clock.Current,
			*scan.VerifiedAt,
		)
	}
}

func TestVerifiedScanStoresVerificationMetadata(t *testing.T) {

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

	// CUSTOMER
	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY WITHOUT VERIFICATION
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
	ticket.UserID = customer.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN
	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// ASSERT STATUS

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified status, got %s",
			scan.Status,
		)
	}

	// ASSERT VERIFIED TIMESTAMP

	if scan.VerifiedAt == nil {

		t.Fatal(
			"expected VerifiedAt to be set",
		)
	}

	if !scan.VerifiedAt.Equal(clock.Current) {

		t.Fatalf(
			"expected VerifiedAt %v, got %v",
			clock.Current,
			*scan.VerifiedAt,
		)
	}

	// ASSERT VERIFY USER

	if scan.VerifiedByID == nil {

		t.Fatal(
			"expected VerifiedByID to be set",
		)
	}

	if *scan.VerifiedByID != staff.ID {

		t.Fatalf(
			"expected VerifiedByID %s, got %s",
			staff.ID,
			*scan.VerifiedByID,
		)
	}

	// ASSERT EXPIRY

	if scan.VerifiedUntil == nil {

		t.Fatal(
			"expected VerifiedUntil to be set",
		)
	}

	expectedExpiry := clock.Current.Add(
		15 * time.Minute,
	)

	if !scan.VerifiedUntil.Equal(expectedExpiry) {

		t.Fatalf(
			"expected VerifiedUntil %v, got %v",
			expectedExpiry,
			*scan.VerifiedUntil,
		)
	}
}
