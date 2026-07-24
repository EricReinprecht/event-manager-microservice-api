package verification

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

func TestVerificationFailsExactlyAtExpiryBoundary(t *testing.T) {

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

		Role: enum.RoleStaff,
	}

	if err := db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET CATEGORY WITH VERIFICATION

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

	if scan.VerificationExpiresAt == nil {
		t.Fatal(
			"expected verification expiry to be set",
		)
	}

	// MOVE EXACTLY TO EXPIRY TIME

	var storedScan appModels.TicketScan

	if err := db.First(
		&storedScan,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	clock.Current = *storedScan.VerificationExpiresAt

	// VERIFY AT EXACT BOUNDARY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err != nil {

		t.Fatalf(
			"expected verification to succeed at expiry boundary, got %v",
			err,
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
			"expected verified status, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedAt == nil {

		t.Fatal(
			"expected VerifiedAt to be set",
		)
	}

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected VerifiedByID to be set",
		)
	}
}

func TestPendingScanExpiresCannotBeVerified(t *testing.T) {

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

	if scan.VerificationExpiresAt == nil {

		t.Fatal(
			"expected verification expiry to be set",
		)
	}

	// MOVE AFTER EXPIRY

	clock.Current = scan.VerificationExpiresAt.Add(
		time.Second,
	)

	// VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err != appErrors.ErrTicketVerificationExpired {

		t.Fatalf(
			"expected verification expired error, got %v",
			err,
		)
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

	// SHOULD STILL BE PENDING

	if updated.Status != enum.TicketScanPending {

		t.Fatalf(
			"expected pending status after expiry, got %s",
			updated.Status,
		)
	}

	if updated.VerifiedAt != nil {

		t.Fatal(
			"expected VerifiedAt to remain nil",
		)
	}

	if updated.VerifiedByID != nil {

		t.Fatal(
			"expected VerifiedByID to remain nil",
		)
	}
}

func TestVerificationExpiredPendingScanCannotBeVerified(t *testing.T) {

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

	// MOVE TIME FORWARD BEYOND VERIFICATION TTL

	clock.Current = clock.Current.Add(
		30 * time.Minute,
	)

	// TRY VERIFY

	err = ticketService.VerifyScan(
		context.Background(),
		scan.ID,
		staff.ID,
		true,
	)

	if err == nil {

		t.Fatal(
			"expected expired verification error",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrTicketVerificationExpired,
	) {

		t.Fatalf(
			"expected ErrTicketVerificationExpired, got %v",
			err,
		)
	}

	// ENSURE STATUS DID NOT CHANGE

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
