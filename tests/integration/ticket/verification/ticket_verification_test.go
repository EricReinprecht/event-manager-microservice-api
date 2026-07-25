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
		rejected,
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

	ticketService := helpers.NewTicketService(db)

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
	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// USERS

	organizer := fixtures.User()

	if err := db.Create(&organizer).Error; err != nil {
		t.Fatal(err)
	}

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	otherUser := fixtures.User()

	if err := db.Create(&otherUser).Error; err != nil {
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

	// PARTY STAFF

	staffMember := appModels.PartyMember{

		ID: uuid.New(),

		UserID: organizer.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&staffMember).Error; err != nil {
		t.Fatal(err)
	}

	staffRole := appModels.PartyMemberRole{

		ID: uuid.New(),

		PartyMemberID: staffMember.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&staffRole).Error; err != nil {
		t.Fatal(err)
	}

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

	// TICKET

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: customer.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
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

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		organizer.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// OTHER PARTY USER TRIES VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		otherUser.ID,
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
		Current: time.Now().UTC(),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// STAFF USER

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

	// MEMBER
	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: customer.ID,
	}

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

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// REJECT SCAN

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// TRY VERIFY AGAIN

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
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

	/// MEMBER
	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: customer.ID,
	}

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

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err != nil {
		t.Fatal(err)
	}

	// TRY REJECT AFTER VERIFIED

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY REQUIRING VERIFICATION

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

	// CREATE PENDING SCAN BY STAFF

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// CUSTOMER TRIES TO VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		customer.ID,
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

	if err := db.First(
		&refreshed,
		" id = ? ",
		scan.ID,
	).Error; err != nil {

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
		Current: time.Now().UTC(),
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

	// ORGANIZER MEMBERSHIP

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: organizer.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	role := appModels.PartyMemberRole{

		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleOrganizer,
	}

	if err := db.Create(&role).Error; err != nil {
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

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		organizer.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// ORGANIZER VERIFIES

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		organizer.ID,
		true,
	)

	if err != nil {
		t.Fatal(err)
	}

	// CHECK RESULT

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected VERIFIED, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected verified_by_id to be set",
		)
	}

	if *updated.VerifiedByID != organizer.ID {

		t.Fatalf(
			"expected verified by organizer %s, got %s",
			organizer.ID,
			*updated.VerifiedByID,
		)
	}

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected verified_at to be set",
		)
	}
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
		Current: time.Now().UTC(),
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// ADMIN

	admin := fixtures.User()

	if err := db.Create(&admin).Error; err != nil {
		t.Fatal(err)
	}

	// CUSTOMER

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY REQUIRED BY PARTY

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

	// ADMIN MEMBERSHIP

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: admin.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	role := appModels.PartyMemberRole{

		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleAdmin,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		admin.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// ADMIN VERIFIES

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		admin.ID,
		true,
	)

	if err != nil {
		t.Fatal(err)
	}

	// VERIFY RESULT

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected VERIFIED, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected verified_by_id",
		)
	}

	if *updated.VerifiedByID != admin.ID {

		t.Fatalf(
			"expected verified by admin %s, got %s",
			admin.ID,
			*updated.VerifiedByID,
		)
	}

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected verified_at",
		)
	}
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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// REJECT

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// LOAD UPDATED SCAN

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	// STATUS

	if updated.Status != enum.TicketScanRejected {

		t.Fatalf(
			"expected rejected status, got %s",
			updated.Status,
		)
	}

	// VERIFIED AT

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected VerifiedAt to be set",
		)
	}

	if !updated.VerifiedAt.UTC().
		Truncate(time.Microsecond).
		Equal(clock.Current.UTC().Truncate(time.Microsecond)) {

		t.Fatalf(
			"expected VerifiedAt %v, got %v",
			clock.Current.UTC(),
			updated.VerifiedAt.UTC(),
		)
	}

	// VERIFIED BY

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected VerifiedByID to be set",
		)
	}

	if *updated.VerifiedByID != staff.ID {

		t.Fatalf(
			"expected VerifiedByID %s, got %s",
			staff.ID,
			*updated.VerifiedByID,
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

	// MEMBER
	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

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
		staff.ID,
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

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	// STAFF WHO CREATED THE SCAN

	staffA := fixtures.User()

	if err := db.Create(&staffA).Error; err != nil {
		t.Fatal(err)
	}

	// OTHER STAFF FROM DIFFERENT PARTY

	staffB := fixtures.User()

	if err := db.Create(&staffB).Error; err != nil {
		t.Fatal(err)
	}

	// CUSTOMER

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY A

	categoryA := fixtures.Category()

	if err := db.Create(&categoryA).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY A

	partyA := fixtures.PartyWithOrganizer(
		staffA.ID,
	)

	partyA.CategoryID = categoryA.ID

	if err := db.Create(&partyA).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER A

	memberA := appModels.PartyMember{

		ID: uuid.New(),

		UserID: staffA.ID,

		PartyID: partyA.ID,
	}

	if err := db.Create(&memberA).Error; err != nil {
		t.Fatal(err)
	}

	roleA := appModels.PartyMemberRole{

		ID: uuid.New(),

		PartyMemberID: memberA.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&roleA).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY B (different UUID!)

	categoryB := appModels.Category{

		ID: uuid.New(),

		Name: "Another Festival",
	}

	if err := db.Create(&categoryB).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY B

	partyB := fixtures.PartyWithOrganizer(
		staffB.ID,
	)

	partyB.CategoryID = categoryB.ID

	if err := db.Create(&partyB).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBER B

	memberB := appModels.PartyMember{

		ID: uuid.New(),

		UserID: staffB.ID,

		PartyID: partyB.ID,
	}

	if err := db.Create(&memberB).Error; err != nil {
		t.Fatal(err)
	}

	roleB := appModels.PartyMemberRole{

		ID: uuid.New(),

		PartyMemberID: memberB.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&roleB).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: partyA.ID,

		RequiresVerification: true,
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

	// STAFF A SCANS

	scan, err := ticketService.Scan(
		context.Background(),
		staffA.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// STAFF B TRIES TO VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staffB.ID,
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

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected scan to remain pending, got %s",
			updated.Status,
		)
	}
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

	// MEMBER
	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	// CREATE PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// REJECT

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// LOAD UPDATED

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	// STATUS

	if updated.Status != enum.TicketScanRejected {

		t.Fatalf(
			"expected rejected status, got %s",
			updated.Status,
		)
	}

	// VERIFICATION METADATA MUST BE CLEARED

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected VerifiedByID to store rejecting staff",
		)
	}

	if *updated.VerifiedByID != staff.ID {

		t.Fatalf(
			"expected VerifiedByID %s, got %s",
			staff.ID,
			*updated.VerifiedByID,
		)
	}

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

	// MEMBER
	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	// WINDOW ONE

	windowOne := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&windowOne).Error; err != nil {
		t.Fatal(err)
	}

	// WINDOW TWO

	windowTwo := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(2 * time.Hour),

		EndsAt: clock.Current.Add(3 * time.Hour),
	}

	if err := db.Create(&windowTwo).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = customer.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// SCAN IN WINDOW ONE

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.TicketAccessWindowID != windowOne.ID {

		t.Fatal(
			"expected scan in first window",
		)
	}

	// VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err != nil {
		t.Fatal(err)
	}

	// MOVE TO SECOND WINDOW

	clock.Current = windowTwo.StartsAt

	// TRY REJECTING OLD VERIFIED SCAN

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		false,
	)

	if err == nil {

		t.Fatal(
			"expected rejection of verified scan to fail",
		)
	}

	// LOAD UPDATED

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified status to remain unchanged, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected VerifiedAt to remain set",
		)
	}

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected VerifiedByID to remain set",
		)
	}
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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY WITHOUT VERIFICATION

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

	// CUSTOMER

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = customer.ID

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

	if scan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified scan, got %s",
			scan.Status,
		)
	}

	if scan.VerifiedAt == nil {

		t.Fatal(
			"expected VerifiedAt to be set",
		)
	}

	// CHANGE CATEGORY SETTING AFTER SCAN

	if err := db.Model(
		&ticketCategory,
	).Update(
		"requires_verification",
		true,
	).Error; err != nil {

		t.Fatal(err)
	}

	// RELOAD SCAN

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	// VERIFY SNAPSHOT DID NOT CHANGE

	if updated.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected existing scan to remain verified, got %s",
			updated.Status,
		)
	}

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
		Current: time.Now().UTC(),
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

	// ORGANIZER MEMBERSHIP
	organizerMember := appModels.PartyMember{

		ID: uuid.New(),

		UserID: organizer.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&organizerMember).Error; err != nil {
		t.Fatal(err)
	}

	organizerRole := appModels.PartyMemberRole{

		ID: uuid.New(),

		PartyMemberID: organizerMember.ID,

		Role: enum.RoleOrganizer,
	}

	if err := db.Create(&organizerRole).Error; err != nil {
		t.Fatal(err)
	}

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

	// CREATE PENDING SCAN AS ORGANIZER

	scan, err := ticketService.Scan(
		context.Background(),
		organizer.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {
		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// ORGANIZER REJECTS

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		organizer.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// VERIFY RESULT

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {
		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanRejected {

		t.Fatalf(
			"expected rejected status, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected VerifiedAt to be set after rejection",
		)
	}

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected VerifiedByID to store rejecting organizer",
		)
	}

	if *updated.VerifiedByID != organizer.ID {

		t.Fatal(
			"expected rejecting organizer to be stored",
		)
	}

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

func TestVerifyAfterAccessWindowClosed(t *testing.T) {

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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Minute),
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

	// CREATE PENDING SCAN INSIDE WINDOW

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

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

	// CHECK RESULT

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified status, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected verification timestamp",
		)
	}

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected verifier metadata",
		)
	}
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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBER ROLE

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

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

	// STAFF CREATES PENDING SCAN

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	// REMOVE STAFF MEMBERSHIP

	if err := db.Delete(
		&member,
	).Error; err != nil {

		t.Fatal(err)
	}

	// TRY VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err != appErrors.ErrNotAllowed {

		t.Fatalf(
			"expected ErrNotAllowed after membership deletion, got %v",
			err,
		)
	}

	// VERIFY SCAN WAS NOT CHANGED

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending status, got %s",
			updated.Status,
		)
	}

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
