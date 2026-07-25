package lifecycle

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

func TestVerifiedTicketCannotBeScannedTwiceInSameAccessWindow(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(db)

	// User
	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	// Category
	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// Party
	party := fixtures.PartyWithOrganizer(
		staff.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// Party Member
	member := appModels.PartyMember{
		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	// Ticket category
	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// Access window
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

	// Purchase
	purchase := appModels.Purchase{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// Ticket
	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,

		PurchaseID: purchase.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// First scan
	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified first scan, got %s",
			firstScan.Status,
		)
	}

	// Second scan should fail
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
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}
}

func TestPendingTicketCannotBeScannedTwiceInSameAccessWindow(t *testing.T) {

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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	role := appModels.PartyMemberRole{
		ID: uuid.New(),

		PartyMemberID: member.ID,

		Role: enum.RoleStaff,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: true,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

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

	// Purchase
	purchase := appModels.Purchase{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// Ticket
	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,

		PurchaseID: purchase.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending first scan, got %s",
			firstScan.Status,
		)
	}

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected duplicate pending scan error",
		)
	}

	if err != appErrors.ErrTicketAlreadyScanned {

		t.Fatalf(
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}

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
			"expected exactly one ticket scan, got %d",
			count,
		)
	}
}

func TestRejectedTicketCanBeScannedAgain(t *testing.T) {

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
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

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

	// Purchase
	purchase := appModels.Purchase{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// Ticket
	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,

		PurchaseID: purchase.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// FIRST SCAN

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected first scan pending, got %s",
			firstScan.Status,
		)
	}

	// REJECT FIRST SCAN

	err = ticketService.VerifyScan(
		context.Background(),
		firstScan.ID,
		staff.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// SECOND SCAN

	secondScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.ID == firstScan.ID {

		t.Fatal(
			"expected a new scan",
		)
	}

	if secondScan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			secondScan.Status,
		)
	}

}

func TestRejectedScanCreatesHistory(t *testing.T) {

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

	// STAFF MEMBER

	member := appModels.PartyMember{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

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

	// PURCHASE

	purchase := appModels.Purchase{

		ID: uuid.New(),

		UserID: customer.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

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

	// FIRST SCAN

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending scan, got %s",
			firstScan.Status,
		)
	}

	// REJECT

	err = ticketService.VerifyScan(
		context.Background(),
		firstScan.ID,
		staff.ID,
		false,
	)

	if err != nil {
		t.Fatal(err)
	}

	// SECOND SCAN AFTER REJECTION

	secondScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.ID == firstScan.ID {

		t.Fatal(
			"expected new scan record",
		)
	}

	// CHECK HISTORY

	var scans []appModels.TicketScan

	err = db.
		Where(
			"ticket_id = ?",
			ticket.ID,
		).
		Order(
			"created_at ASC",
		).
		Find(&scans).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(scans) != 2 {

		t.Fatalf(
			"expected 2 scan records, got %d",
			len(scans),
		)
	}

	if scans[0].Status != enum.TicketScanRejected {

		t.Fatalf(
			"expected first scan rejected, got %s",
			scans[0].Status,
		)
	}

	if scans[1].Status != enum.TicketScanPending {

		t.Fatalf(
			"expected second scan pending, got %s",
			scans[1].Status,
		)
	}
}

func TestStaffCannotScanTicketAfterVerificationExpired(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	fakeClock := &helpers.FakeClock{
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

	member := appModels.PartyMember{}

	err = db.
		Where(
			"user_id = ? AND party_id = ?",
			staff.ID,
			party.ID,
		).
		First(&member).
		Error

	if err != nil {

		// no member exists yet, create one

		member = appModels.PartyMember{

			ID: uuid.New(),

			UserID: staff.ID,

			PartyID: party.ID,
		}

		if err := db.Create(&member).Error; err != nil {
			t.Fatal(err)
		}
	}

	// Add STAFF role

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

		StartsAt: fakeClock.Current.Add(-time.Hour),

		EndsAt: fakeClock.Current.Add(time.Hour * 2),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	// Purchase
	purchase := appModels.Purchase{

		ID: uuid.New(),

		UserID: staff.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// Ticket
	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,

		PurchaseID: purchase.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// FIRST SCAN

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected verified scan, got %s",
			firstScan.Status,
		)
	}

	// SECOND SCAN BEFORE EXPIRY

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected duplicate scan before expiry",
		)
	}

	if err != appErrors.ErrTicketAlreadyScanned {

		t.Fatalf(
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}

	// MOVE CLOCK AFTER VERIFICATION TTL

	fakeClock.Current = fakeClock.Current.Add(
		16 * time.Minute,
	)

	// THIRD SCAN AFTER EXPIRY

	secondScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.Status != enum.TicketScanVerified {

		t.Fatalf(
			"expected new verified scan after expiry, got %s",
			secondScan.Status,
		)
	}

}
