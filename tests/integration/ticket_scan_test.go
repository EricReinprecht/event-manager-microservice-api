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

func TestStaffCanScanTicketWithoutVerification(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

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

		Price: 100,

		Capacity: 100,

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

		UserID: staff.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN

	result, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.TicketID != ticket.ID {
		t.Fatalf(
			"wrong ticket scanned. expected %s got %s",
			ticket.ID,
			result.TicketID,
		)
	}

	// VERIFY DATABASE ENTRY

	var scan appModels.TicketScan

	err = db.
		Where(
			"ticket_id = ?",
			ticket.ID,
		).
		First(&scan).
		Error

	if err != nil {

		t.Fatal("ticket scan was not created")
	}

	if scan.Status != enum.TicketScanVerified {
		t.Fatalf(
			"expected verified status, got %s",
			scan.Status,
		)
	}

	if scan.VerifiedAt == nil {

		t.Fatal(
			"expected verified_at to be set",
		)
	}

	if scan.VerifiedByID == nil {

		t.Fatal(
			"expected verified_by_id to be set",
		)
	}

	if *scan.VerifiedByID != staff.ID {

		t.Fatalf(
			"wrong verifier. expected %s got %s",
			staff.ID,
			*scan.VerifiedByID,
		)
	}

}

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

func TestTwoDifferentTicketsCanBeScannedBySameStaff(t *testing.T) {

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

	// CUSTOMERS

	customerOne := fixtures.User()

	if err := db.Create(&customerOne).Error; err != nil {
		t.Fatal(err)
	}

	customerTwo := fixtures.User()

	if err := db.Create(&customerTwo).Error; err != nil {
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

	// STAFF MEMBER

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

	// TICKET ONE

	ticketOne := fixtures.Ticket()

	ticketOne.TicketCategoryID = ticketCategory.ID
	ticketOne.UserID = customerOne.ID

	if err := db.Create(&ticketOne).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET TWO

	ticketTwo := fixtures.Ticket()

	ticketTwo.TicketCategoryID = ticketCategory.ID
	ticketTwo.UserID = customerTwo.ID

	if err := db.Create(&ticketTwo).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN FIRST TICKET

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticketOne.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.TicketID != ticketOne.ID {

		t.Fatalf(
			"expected first scan ticket %s, got %s",
			ticketOne.ID,
			firstScan.TicketID,
		)
	}

	// SCAN SECOND TICKET

	secondScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticketTwo.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.TicketID != ticketTwo.ID {

		t.Fatalf(
			"expected second scan ticket %s, got %s",
			ticketTwo.ID,
			secondScan.TicketID,
		)
	}

	// BOTH MUST BE SAME STAFF

	if firstScan.ScannedByID != staff.ID {

		t.Fatalf(
			"expected first scan by staff",
		)
	}

	if secondScan.ScannedByID != staff.ID {

		t.Fatalf(
			"expected second scan by staff",
		)
	}

	// BOTH MUST BE SAME WINDOW

	if firstScan.TicketAccessWindowID != window.ID {

		t.Fatalf(
			"expected first scan in window",
		)
	}

	if secondScan.TicketAccessWindowID != window.ID {

		t.Fatalf(
			"expected second scan in window",
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
