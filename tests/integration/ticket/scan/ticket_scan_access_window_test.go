package scan

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
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func setupAccessWindowTest(
	t *testing.T,
	clock helpers.Clock,
	start time.Time,
	end time.Time,
) (
	*service.TicketService,
	*gorm.DB,
	appModels.Ticket,
	appModels.User,
) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

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

		Price: 100,

		Capacity: 100,

		PartyID: party.ID,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	window := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: start,

		EndsAt: end,
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	ticket := appModels.Ticket{

		ID: uuid.New(),

		Code: uuid.NewString(),

		TicketCategoryID: ticketCategory.ID,

		UserID: staff.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ticketService, db, ticket, staff
}

func TestStaffCannotScanTicketBeforeAccessWindow(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	ticketService, _, ticket, staff := setupAccessWindowTest(
		t,
		fakeClock,
		fakeClock.Now().Add(time.Hour),
		fakeClock.Now().Add(2*time.Hour),
	)

	_, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected scan to fail before access window",
		)
	}

	if !errors.Is(err, appErrors.ErrTicketNotValidNow) {

		t.Fatalf(
			"expected ErrTicketNotValidNow, got %v",
			err,
		)
	}
}

func TestStaffCanScanTicketExactlyAtAccessWindowStart(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	start := fakeClock.Now()

	ticketService, _, ticket, staff := setupAccessWindowTest(
		t,
		fakeClock,
		start,
		start.Add(time.Hour),
	)

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {

		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s got %s",
			ticket.ID,
			scan.TicketID,
		)
	}
}

func TestStaffCanScanTicketExactlyAtAccessWindowEnd(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	end := fakeClock.Now().Add(time.Hour)

	ticketService, _, ticket, staff := setupAccessWindowTest(
		t,
		fakeClock,
		end.Add(-time.Hour),
		end,
	)

	fakeClock.Set(end)

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {

		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s got %s",
			ticket.ID,
			scan.TicketID,
		)
	}
}

func TestStaffCannotScanTicketAfterAccessWindow(t *testing.T) {

	fakeClock := helpers.NewFakeClock(
		time.Date(
			2026,
			7,
			24,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	)

	ticketService, _, ticket, staff := setupAccessWindowTest(
		t,
		fakeClock,
		fakeClock.Now().Add(-2*time.Hour),
		fakeClock.Now().Add(-time.Hour),
	)

	_, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected scan to fail after access window",
		)
	}

	if !errors.Is(err, appErrors.ErrTicketNotValidNow) {

		t.Fatalf(
			"expected ErrTicketNotValidNow, got %v",
			err,
		)
	}
}

func TestScanUsesCorrectAccessWindowWhenWindowsOverlap(t *testing.T) {

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

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// OVERLAPPING WINDOW 1 (starts earlier)

	firstWindow := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-2 * time.Hour),

		EndsAt: clock.Current.Add(2 * time.Hour),
	}

	if err := db.Create(&firstWindow).Error; err != nil {
		t.Fatal(err)
	}

	// OVERLAPPING WINDOW 2 (starts later)

	secondWindow := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&secondWindow).Error; err != nil {
		t.Fatal(err)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = staff.ID

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

	// ASSERT CORRECT WINDOW

	if scan.TicketAccessWindowID != firstWindow.ID {

		t.Fatalf(
			"expected scan to use first access window %s, got %s",
			firstWindow.ID,
			scan.TicketAccessWindowID,
		)
	}
}

func TestScanTicketWithMultipleActiveWindows(t *testing.T) {

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

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACTIVE WINDOW ONE

	windowOne := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-2 * time.Hour),

		EndsAt: clock.Current.Add(2 * time.Hour),
	}

	if err := db.Create(&windowOne).Error; err != nil {
		t.Fatal(err)
	}

	// ACTIVE WINDOW TWO (OVERLAPS)

	windowTwo := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&windowTwo).Error; err != nil {
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

	scan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if scan.TicketID != ticket.ID {

		t.Fatalf(
			"expected ticket %s got %s",
			ticket.ID,
			scan.TicketID,
		)
	}

	// VERIFY ONE WINDOW WAS SELECTED

	if scan.TicketAccessWindowID != windowOne.ID &&
		scan.TicketAccessWindowID != windowTwo.ID {

		t.Fatalf(
			"expected scan to belong to one active window, got %s",
			scan.TicketAccessWindowID,
		)
	}

	// VERIFY ONLY ONE SCAN EXISTS

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
			"expected exactly one scan, got %d",
			count,
		)
	}
}

func TestScanTicketWithoutAccessWindow(t *testing.T) {

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

	// TICKET CATEGORY WITHOUT ACCESS WINDOWS

	ticketCategory := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// IMPORTANT:
	// No TicketAccessWindow is created.

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

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected scan to fail without access window",
		)
	}

	if !errors.Is(
		err,
		appErrors.ErrTicketNotValidNow,
	) {

		t.Fatalf(
			"expected ErrTicketNotValidNow, got %v",
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
			"expected no ticket scans, got %d",
			count,
		)
	}
}

func TestTicketCanBeScannedInDifferentAccessWindow(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	fakeClock := &helpers.FakeClock{
		Current: time.Now().UTC(),
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

	// FIRST ACCESS WINDOW

	windowOne := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: fakeClock.Current.Add(-time.Hour),

		EndsAt: fakeClock.Current.Add(time.Hour),
	}

	if err := db.Create(&windowOne).Error; err != nil {
		t.Fatal(err)
	}

	// SECOND ACCESS WINDOW

	windowTwo := appModels.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: fakeClock.Current.Add(2 * time.Hour),

		EndsAt: fakeClock.Current.Add(3 * time.Hour),
	}

	if err := db.Create(&windowTwo).Error; err != nil {
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

	// FIRST SCAN IN WINDOW ONE

	firstScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if firstScan.TicketAccessWindowID != windowOne.ID {

		t.Fatalf(
			"expected first scan in window one",
		)
	}

	// SECOND SCAN IN SAME WINDOW SHOULD FAIL

	_, err = ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err == nil {

		t.Fatal(
			"expected duplicate scan in same window",
		)
	}

	// MOVE CLOCK INTO SECOND WINDOW

	fakeClock.Current = fakeClock.Current.Add(
		2 * time.Hour,
	)

	// SCAN AGAIN IN SECOND WINDOW

	secondScan, err := ticketService.Scan(
		context.Background(),
		staff.ID,
		ticket.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if secondScan.TicketAccessWindowID != windowTwo.ID {

		t.Fatalf(
			"expected second scan in window two",
		)
	}

	if secondScan.ID == firstScan.ID {

		t.Fatal(
			"expected a new scan record",
		)
	}
}
