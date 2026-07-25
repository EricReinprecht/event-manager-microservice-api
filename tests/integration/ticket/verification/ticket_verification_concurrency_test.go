package verification

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestConcurrentVerifyAndReject(t *testing.T) {

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

	//  Ticket
	purchase := helpers.CreateTestPurchase(
		db,
		customer.ID,
		party.ID,
	)

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

	if scan.Status != enum.TicketScanPending {
		t.Fatalf(
			"expected pending scan, got %s",
			scan.Status,
		)
	}

	var wg sync.WaitGroup

	results := make(chan error, 2)

	wg.Add(2)

	// VERIFY

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff.ID,
			true,
		)

	}()

	// REJECT

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff.ID,
			false,
		)

	}()

	wg.Wait()

	close(results)

	success := 0
	failures := 0

	for err := range results {

		if err == nil {
			success++
		} else {
			failures++
		}
	}

	if success != 1 {

		t.Fatalf(
			"expected exactly one successful decision, got %d",
			success,
		)
	}

	if failures != 1 {

		t.Fatalf(
			"expected exactly one failed decision, got %d",
			failures,
		)
	}

	// VERIFY FINAL DATABASE STATE

	var updated appModels.TicketScan

	if err := db.First(
		&updated,
		"id = ?",
		scan.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	switch updated.Status {

	case enum.TicketScanVerified:

		if updated.VerifiedAt == nil {
			t.Fatal(
				"verified scan missing decision timestamp",
			)
		}

		if updated.VerifiedByID == nil {
			t.Fatal(
				"verified scan missing decision user",
			)
		}

		if updated.VerifiedUntil == nil {
			t.Fatal(
				"verified scan missing verification expiry",
			)
		}

	case enum.TicketScanRejected:

		if updated.VerifiedAt == nil {
			t.Fatal(
				"rejected scan missing decision timestamp",
			)
		}

		if updated.VerifiedByID == nil {
			t.Fatal(
				"rejected scan missing rejecting user",
			)
		}

		if updated.VerifiedUntil != nil {
			t.Fatal(
				"rejected scan should not have verification validity",
			)
		}

	default:

		t.Fatalf(
			"unexpected final status %s",
			updated.Status,
		)
	}
}

func TestConcurrentVerifySameScan(t *testing.T) {

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

	// STAFF 1

	staff1 := fixtures.User()

	if err := db.Create(&staff1).Error; err != nil {
		t.Fatal(err)
	}

	// STAFF 2

	staff2 := fixtures.User()

	if err := db.Create(&staff2).Error; err != nil {
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
		staff1.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// MEMBERS

	members := []appModels.PartyMember{

		{
			ID: uuid.New(),

			UserID: staff1.ID,

			PartyID: party.ID,
		},

		{
			ID: uuid.New(),

			UserID: staff2.ID,

			PartyID: party.ID,
		},
	}

	for _, member := range members {

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

	//  Purchase
	purchase := appModels.Purchase{

		ID: uuid.New(),

		UserID: customer.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// Purchase
	purchase = appModels.Purchase{

		ID: uuid.New(),

		UserID: customer.ID,

		PartyID: party.ID,

		Status: enum.PurchaseStatusPaid,
	}

	if err := db.Create(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	//  Ticket
	purchase = helpers.CreateTestPurchase(
		db,
		customer.ID,
		party.ID,
	)

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
		staff1.ID,
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

	// RUN TWO APPROVALS CONCURRENTLY

	results := make(chan error, 2)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff1.ID,
			true,
		)

	}()

	go func() {

		defer wg.Done()

		results <- ticketService.VerifyScan(
			context.Background(),
			scan.ID,
			staff2.ID,
			true,
		)

	}()

	wg.Wait()

	close(results)

	successes := 0
	failures := 0

	for err := range results {

		if err == nil {

			successes++

		} else if err == appErrors.ErrTicketScanAlreadyDecided {

			failures++

		} else {

			t.Fatalf(
				"unexpected error: %v",
				err,
			)
		}
	}

	if successes != 1 {

		t.Fatalf(
			"expected exactly one successful verification, got %d",
			successes,
		)
	}

	if failures != 1 {

		t.Fatalf(
			"expected exactly one rejected verification attempt, got %d",
			failures,
		)
	}

	// VERIFY FINAL STATE

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

	if updated.VerifiedByID == nil {

		t.Fatal(
			"expected verifying staff metadata",
		)
	}
}
