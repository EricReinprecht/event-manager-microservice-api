package scan

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

func TestTwoStaffCannotScanSameTicketAtSameTime(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

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

	// USERS

	staffOne := fixtures.User()

	if err := db.Create(&staffOne).Error; err != nil {
		t.Fatal(err)
	}

	staffTwo := fixtures.User()

	if err := db.Create(&staffTwo).Error; err != nil {
		t.Fatal(err)
	}

	// CATEGORY

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY

	party := fixtures.PartyWithOrganizer(
		staffOne.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBERS

	members := []appModels.PartyMember{

		{
			ID: uuid.New(),

			UserID: staffOne.ID,

			PartyID: party.ID,

			Role: enum.RoleStaff,
		},

		{
			ID: uuid.New(),

			UserID: staffTwo.ID,

			PartyID: party.ID,

			Role: enum.RoleStaff,
		},
	}

	for _, member := range members {

		if err := db.Create(&member).Error; err != nil {
			t.Fatal(err)
		}
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

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staffOne.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	// RUN CONCURRENT SCANS

	var wg sync.WaitGroup

	wg.Add(2)

	results := make(chan error, 2)

	go func() {

		defer wg.Done()

		_, err := ticketService.Scan(
			context.Background(),
			staffOne.ID,
			ticket.Code,
		)

		results <- err
	}()

	go func() {

		defer wg.Done()

		_, err := ticketService.Scan(
			context.Background(),
			staffTwo.ID,
			ticket.Code,
		)

		results <- err
	}()

	wg.Wait()

	close(results)

	successes := 0
	duplicates := 0

	for err := range results {

		if err == nil {

			successes++

			continue
		}

		if err == appErrors.ErrTicketAlreadyScanned {

			duplicates++

			continue
		}

		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if successes != 1 {

		t.Fatalf(
			"expected exactly one successful scan, got %d",
			successes,
		)
	}

	if duplicates != 1 {

		t.Fatalf(
			"expected exactly one duplicate error, got %d",
			duplicates,
		)
	}

	// DATABASE CHECK

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
			"expected one scan in database, got %d",
			count,
		)
	}
}

func TestSameTicketDifferentStaffSameWindow(t *testing.T) {

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

	// STAFF ONE

	staffOne := fixtures.User()

	if err := db.Create(&staffOne).Error; err != nil {
		t.Fatal(err)
	}

	// STAFF TWO

	staffTwo := fixtures.User()

	if err := db.Create(&staffTwo).Error; err != nil {
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
		staffOne.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	// PARTY MEMBERS

	members := []appModels.PartyMember{

		{
			ID: uuid.New(),

			UserID: staffOne.ID,

			PartyID: party.ID,

			Role: enum.RoleStaff,
		},

		{
			ID: uuid.New(),

			UserID: staffTwo.ID,

			PartyID: party.ID,

			Role: enum.RoleStaff,
		},
	}

	for _, member := range members {

		if err := db.Create(&member).Error; err != nil {
			t.Fatal(err)
		}
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

	// FIRST STAFF SCAN

	firstScan, err := ticketService.Scan(
		context.Background(),
		staffOne.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	// SECOND STAFF SCAN

	secondScan, err := ticketService.Scan(
		context.Background(),
		staffTwo.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected second staff scan to fail",
		)
	}

	if err != appErrors.ErrTicketAlreadyScanned {

		t.Fatalf(
			"expected ErrTicketAlreadyScanned, got %v",
			err,
		)
	}

	// FIRST SCAN MUST BE KEPT

	if firstScan.TicketID != ticket.ID {

		t.Fatalf(
			"expected scan for ticket %s",
			ticket.ID,
		)
	}

	if firstScan.TicketAccessWindowID != window.ID {

		t.Fatalf(
			"expected scan in current window",
		)
	}

	// VERIFY ONLY ONE SCAN EXISTS

	var count int64

	if err := db.Model(
		&appModels.TicketScan{},
	).
		Where(
			"ticket_id = ? AND ticket_access_window_id = ?",
			ticket.ID,
			window.ID,
		).
		Count(&count).Error; err != nil {

		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected exactly one scan, got %d",
			count,
		)
	}

	_ = secondScan
}
