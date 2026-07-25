package verification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestStaffCanRejectPendingTicketScan(t *testing.T) {

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

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	// ASSERT INITIAL STATE

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// REJECT

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		false,
	)

	// RELOAD

	var rejected appModels.TicketScan

	if err := db.
		First(
			&rejected,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	assertRejected(
		t,
		&rejected,
		s.Staff.ID,
	)
}

func TestStaffCanVerifyPendingTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			25,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	// SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// VERIFY

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		true,
	)

	// CHECK RESULT

	var verified appModels.TicketScan

	if err := db.
		First(
			&verified,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	assertVerified(
		t,
		&verified,
		s.Staff.ID,
	)
}

func TestVerificationCanOnlyBeDoneByPartyStaff(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN
	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// OTHER USER TRIES VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		s.OtherUser.ID,
		true,
	)

	if err == nil {

		t.Fatal(
			"expected verification permission error",
		)
	}

	if err != appErrors.ErrNotAllowed {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}
}

func TestCannotVerifyAlreadyRejectedScan(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// REJECT SCAN

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		false,
	)

	// TRY VERIFY AGAIN

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		s.Staff.ID,
		true,
	)

	if err == nil {

		t.Fatal(
			"expected already decided scan error",
		)
	}

	if err != appErrors.ErrTicketScanAlreadyDecided {

		t.Fatalf(
			"expected ErrTicketScanAlreadyDecided, got %v",
			err,
		)
	}
}

func TestCannotRejectAlreadyVerifiedScan(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// VERIFY SCAN

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		true,
	)

	// TRY REJECT AFTER VERIFIED

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		s.Staff.ID,
		false,
	)

	if err == nil {

		t.Fatal(
			"expected already decided scan error",
		)
	}

	if err != appErrors.ErrTicketScanAlreadyDecided {

		t.Fatalf(
			"expected ErrTicketScanAlreadyDecided, got %v",
			err,
		)
	}
}

func TestCustomerCannotVerifyPendingTicketScan(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN BY STAFF

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// CUSTOMER TRIES TO VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		s.Customer.ID,
		true,
	)

	if err == nil {

		t.Fatal(
			"expected customer verification to fail",
		)
	}

	if err != appErrors.ErrNotAllowed {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}

	// VERIFY SCAN DID NOT CHANGE

	var refreshed appModels.TicketScan

	if err := db.
		First(
			&refreshed,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	if refreshed.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected scan to remain pending, got %s",
			refreshed.Status,
		)
	}
}

func TestOrganizerCanVerifyPendingTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	makeOrganizer(
		t,
		db,
		s,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// ORGANIZER VERIFIES

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		true,
	)

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	assertVerified(
		t,
		&updated,
		s.Staff.ID,
	)
}

func TestAdminCanVerifyPendingTicket(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	// CHANGE STAFF ROLE TO ADMIN

	makeAdmin(
		t,
		db,
		s,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// ADMIN VERIFY

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		true,
	)

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	assertVerified(
		t,
		&updated,
		s.Staff.ID,
	)
}

func TestRejectedScanStoresRejectionMetadata(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// REJECT

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		false,
	)

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// BASIC ASSERTIONS

	assertRejected(
		t,
		&updated,
		s.Staff.ID,
	)

	// VERIFIED AT

	if !updated.VerifiedAt.UTC().
		Truncate(time.Microsecond).
		Equal(clock.Current.UTC().Truncate(time.Microsecond)) {

		t.Fatalf(
			"expected VerifiedAt %v, got %v",
			clock.Current.UTC(),
			updated.VerifiedAt.UTC(),
		)
	}

	// EXPIRY CLEARED

	if updated.VerificationExpiresAt != nil {

		t.Fatal(
			"expected VerificationExpiresAt to be nil after rejection",
		)
	}

	// VERIFIED UNTIL CLEARED

	if updated.VerifiedUntil != nil {

		t.Fatal(
			"expected VerifiedUntil to be nil after rejection",
		)
	}
}

func TestDeletedScanCannotBeVerified(t *testing.T) {

	db, err := helpers.TestDatabase()
	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// DELETE SCAN

	if err := db.Delete(
		&appModels.TicketScan{},
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	// VERIFY DELETED SCAN

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		s.Staff.ID,
		true,
	)

	if err == nil {

		t.Fatal(
			"expected deleted scan verification error",
		)
	}

	if !errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {

		t.Fatalf(
			"expected gorm.ErrRecordNotFound, got %v",
			err,
		)
	}

	// ENSURE IT IS STILL DELETED

	var count int64

	if err := db.
		Model(&appModels.TicketScan{}).
		Where(
			"id = ?",
			scan.ID,
		).
		Count(&count).
		Error; err != nil {

		t.Fatal(err)
	}

	if count != 0 {

		t.Fatalf(
			"expected deleted scan to remain deleted, got %d records",
			count,
		)
	}
}

func TestVerificationCannotBeDoneByStaffFromAnotherParty(t *testing.T) {

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

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// SECOND PARTY

	secondParty := createSecondParty(
		t,
		db,
		s.OtherUser,
	)

	// STAFF FROM SECOND PARTY

	secondPartyStaff := addSecondPartyStaff(
		t,
		db,
		secondParty.ID,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// OTHER PARTY STAFF TRIES TO VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		secondPartyStaff.ID,
		true,
	)

	if err == nil {

		t.Fatal(
			"expected verification to fail",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed, got %v",
			err,
		)
	}

	// ENSURE SCAN WAS NOT CHANGED

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	assertScanStatus(
		t,
		&updated,
		enum.TicketScanPending,
	)
}

func TestRejectedScanClearsVerificationExpiryMetadata(t *testing.T) {

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

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// REJECT

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		false,
	)

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	assertRejected(
		t,
		&updated,
		s.Staff.ID,
	)

	// VERIFICATION METADATA MUST BE CLEARED

	if updated.VerificationExpiresAt != nil {

		t.Fatal(
			"expected VerificationExpiresAt to be nil after rejection",
		)
	}

	if updated.VerifiedUntil != nil {

		t.Fatal(
			"expected VerifiedUntil to be nil after rejection",
		)
	}
}

func TestVerifiedScanCannotBeRejectedByDifferentWindow(t *testing.T) {

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

	s := createVerificationScenario(
		t,
		db,
		clock,
	)

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// SECOND WINDOW

	windowTwo := createAccessWindow(
		t,
		db,
		s.TicketCategory.ID,
		clock.Current.Add(2*time.Hour),
		clock.Current.Add(3*time.Hour),
	)

	// CREATE PENDING SCAN

	scan := createPendingScan(
		t,
		ticketService,
		s.Staff.ID,
		s.Ticket,
	)

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	if scan.TicketAccessWindowID != s.Window.ID {

		t.Fatal(
			"expected scan in first window",
		)
	}

	// VERIFY

	verifyScan(
		t,
		ticketService,
		scan,
		s.Staff.ID,
		true,
	)

	// MOVE TO SECOND WINDOW

	clock.Current = windowTwo.StartsAt

	// TRY REJECTING VERIFIED SCAN

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		s.Staff.ID,
		false,
	)

	if err == nil {

		t.Fatal(
			"expected rejection of verified scan to fail",
		)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	assertVerified(
		t,
		&updated,
		s.Staff.ID,
	)
}

func TestTicketCategoryVerificationSettingSnapshot(t *testing.T) {

	db, err := helpers.TestDatabase()
	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
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

	// STAFF ROLE

	addPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

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

	createAccessWindow(
		t,
		db,
		ticketCategory.ID,
		clock.Current.Add(-time.Hour),
		clock.Current.Add(time.Hour),
	)

	// PURCHASE

	purchase := fixtures.Purchase(
		customer.ID,
		party.ID,
	)

	// make sure FK values are populated
	purchase.UserID = customer.ID
	purchase.PartyID = party.ID

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = customer.ID
	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN WHILE VERIFICATION DISABLED

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	assertScanStatus(
		t,
		scan,
		enum.TicketScanVerified,
	)

	if scan.VerifiedAt == nil {
		t.Fatal(
			"expected VerifiedAt to be set",
		)
	}

	// CHANGE CATEGORY SETTING AFTER SCAN

	if err := db.
		Model(&ticketCategory).
		Update(
			"requires_verification",
			true,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// RELOAD SCAN

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// SNAPSHOT ASSERTION

	assertScanStatus(
		t,
		&updated,
		enum.TicketScanVerified,
	)

	if updated.VerifiedAt == nil {
		t.Fatal(
			"expected verification metadata to remain",
		)
	}

	if updated.VerificationExpiresAt != nil {
		t.Fatal(
			"expected no verification expiry for already verified scan",
		)
	}
}

func TestOrganizerCanRejectPendingTicket(t *testing.T) {

	db, err := helpers.TestDatabase()
	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// ORGANIZER

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
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
		organizer.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// ORGANIZER ROLE

	addPartyRole(
		t,
		db,
		organizer.ID,
		party.ID,
		enum.RoleOrganizer,
	)

	// TICKET CATEGORY

	ticketCategory := appModels.TicketCategory{
		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: true,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	createAccessWindow(
		t,
		db,
		ticketCategory.ID,
		clock.Current.Add(-time.Hour),
		clock.Current.Add(time.Hour),
	)

	// PURCHASE

	purchase := fixtures.Purchase(
		customer.ID,
		party.ID,
	)

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = customer.ID
	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		organizer.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// REJECT

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		organizer.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	assertRejected(
		t,
		&updated,
		organizer.ID,
	)
}

func TestVerifyAfterAccessWindowClosed(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
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

	// STAFF ROLE

	addPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	// TICKET CATEGORY

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: true,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	window := createAccessWindow(
		t,
		db,
		ticketCategory.ID,
		clock.Current.Add(-time.Hour),
		clock.Current.Add(time.Minute),
	)

	// PURCHASE

	purchase := fixtures.Purchase(
		customer.ID,
		party.ID,
	)

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = customer.ID
	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// CREATE PENDING SCAN INSIDE WINDOW

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// MOVE AFTER ACCESS WINDOW CLOSED

	clock.Current = window.EndsAt.Add(time.Minute)

	// VERIFY AFTER WINDOW CLOSED

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err != nil {
		t.Fatalf(
			"expected verification after window close to succeed, got %v",
			err,
		)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT

	assertVerified(
		t,
		&updated,
		staff.ID,
	)
}

func TestScannerDeletedBeforeVerification(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	clock := &helpers.FakeClock{
		Current: time.Now().
			UTC().
			Truncate(time.Microsecond),
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

	// STAFF ROLE

	addPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	// TICKET CATEGORY

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: true,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	createAccessWindow(
		t,
		db,
		ticketCategory.ID,
		clock.Current.Add(-time.Hour),
		clock.Current.Add(time.Hour),
	)

	// PURCHASE

	purchase := fixtures.Purchase(
		customer.ID,
		party.ID,
	)

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = customer.ID
	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	assertScanStatus(
		t,
		scan,
		enum.TicketScanPending,
	)

	// REMOVE STAFF MEMBERSHIP

	member := appModels.PartyMember{}

	if err := db.
		Where(
			"user_id = ? AND party_id = ?",
			staff.ID,
			party.ID,
		).
		First(&member).
		Error; err != nil {

		t.Fatal(err)
	}

	// remove roles first because of FK constraint
	if err := db.
		Where(
			"party_member_id = ?",
			member.ID,
		).
		Delete(
			&appModels.PartyMemberRole{},
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// now remove membership
	if err := db.Delete(&member).Error; err != nil {
		t.Fatal(err)
	}

	// TRY VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if !errors.Is(
		err,
		appErrors.ErrNotAllowed,
	) {

		t.Fatalf(
			"expected ErrNotAllowed after membership deletion, got %v",
			err,
		)
	}

	// RELOAD

	var updated appModels.TicketScan

	if err := db.
		First(
			&updated,
			"id = ?",
			scan.ID,
		).
		Error; err != nil {

		t.Fatal(err)
	}

	// ASSERT UNCHANGED

	assertScanStatus(
		t,
		&updated,
		enum.TicketScanPending,
	)

	if updated.VerifiedAt != nil {

		t.Fatal(
			"expected no verification timestamp",
		)
	}

	if updated.VerifiedByID != nil {

		t.Fatal(
			"expected no verifier",
		)
	}
}
