package scenarios

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/service"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

type AccessWindowScenario struct {
	TicketService *service.TicketService

	DB *gorm.DB

	Ticket appModels.Ticket

	Staff appModels.User

	Window appModels.TicketAccessWindow

	TicketCategory appModels.TicketCategory
}

type OverlappingAccessWindowScenario struct {
	TicketService interface {
		Scan(ctx context.Context, userID uuid.UUID, code string) (*appModels.TicketScan, error)
	}

	DB *gorm.DB

	Clock *helpers.FakeClock

	Staff appModels.User

	Ticket appModels.Ticket

	FirstWindow appModels.TicketAccessWindow

	SecondWindow appModels.TicketAccessWindow
}

type MultipleActiveWindowsScenario struct {
	TicketService *service.TicketService

	DB *gorm.DB

	Staff models.User

	Ticket models.Ticket

	WindowOne models.TicketAccessWindow

	WindowTwo models.TicketAccessWindow
}

type NoAccessWindowScenario struct {
	TicketService *service.TicketService

	DB *gorm.DB

	Staff models.User

	Ticket models.Ticket
}

type SequentialAccessWindowsScenario struct {
	TicketService *service.TicketService

	DB *gorm.DB

	Clock *helpers.FakeClock

	Staff models.User

	Ticket models.Ticket

	WindowOne models.TicketAccessWindow

	WindowTwo models.TicketAccessWindow
}

func CreateAccessWindowScenario(
	t *testing.T,
	clock *helpers.FakeClock,
	start time.Time,
	end time.Time,
) AccessWindowScenario {

	t.Helper()

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

	base := CreateTicketScenario(
		t,
		db,
		clock,
		false,
		WithoutAccessWindow(),
	)

	var windows []appModels.TicketAccessWindow

	db.
		Where(
			"ticket_category_id = ?",
			base.TicketCategory.ID,
		).
		Find(&windows)

	window := appModels.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: base.TicketCategory.ID,

		StartsAt: start.UTC(),

		EndsAt: end.UTC(),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	return AccessWindowScenario{

		TicketService: ticketService,

		DB: db,

		Ticket: base.Ticket,

		Staff: base.Staff,

		Window: window,

		TicketCategory: base.TicketCategory,
	}
}

func CreateOverlappingAccessWindowScenario(
	t *testing.T,
	clock *helpers.FakeClock,
) (
	*service.TicketService,
	OverlappingAccessWindowScenario,
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

	AddPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	ticketCategory := appModels.TicketCategory{
		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	firstWindow := appModels.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-2 * time.Hour),

		EndsAt: clock.Current.Add(2 * time.Hour),
	}

	if err := db.Create(&firstWindow).Error; err != nil {
		t.Fatal(err)
	}

	secondWindow := appModels.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&secondWindow).Error; err != nil {
		t.Fatal(err)
	}

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = staff.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ticketService, OverlappingAccessWindowScenario{

		DB: db,

		Clock: clock,

		Staff: staff,

		Ticket: ticket,

		FirstWindow: firstWindow,

		SecondWindow: secondWindow,
	}
}

func CreateMultipleActiveWindowsScenario(
	t *testing.T,
	clock *helpers.FakeClock,
) MultipleActiveWindowsScenario {

	t.Helper()

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

	AddPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	ticketCategory := models.TicketCategory{
		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	windowOne := models.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-2 * time.Hour),

		EndsAt: clock.Current.Add(2 * time.Hour),
	}

	if err := db.Create(&windowOne).Error; err != nil {
		t.Fatal(err)
	}

	windowTwo := models.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&windowTwo).Error; err != nil {
		t.Fatal(err)
	}

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = staff.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return MultipleActiveWindowsScenario{
		TicketService: ticketService,

		DB: db,

		Staff: staff,

		Ticket: ticket,

		WindowOne: windowOne,

		WindowTwo: windowTwo,
	}
}

func CreateNoAccessWindowScenario(
	t *testing.T,
	clock *helpers.FakeClock,
) NoAccessWindowScenario {

	t.Helper()

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

	AddPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	ticketCategory := models.TicketCategory{
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

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = staff.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return NoAccessWindowScenario{
		TicketService: ticketService,

		DB: db,

		Staff: staff,

		Ticket: ticket,
	}
}

func CreateSequentialAccessWindowsScenario(
	t *testing.T,
	clock *helpers.FakeClock,
) SequentialAccessWindowsScenario {

	t.Helper()

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

	AddPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	ticketCategory := models.TicketCategory{
		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: false,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	windowOne := models.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&windowOne).Error; err != nil {
		t.Fatal(err)
	}

	windowTwo := models.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(2 * time.Hour),

		EndsAt: clock.Current.Add(3 * time.Hour),
	}

	if err := db.Create(&windowTwo).Error; err != nil {
		t.Fatal(err)
	}

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID
	ticket.UserID = staff.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return SequentialAccessWindowsScenario{
		TicketService: ticketService,

		DB: db,

		Clock: clock,

		Staff: staff,

		Ticket: ticket,

		WindowOne: windowOne,

		WindowTwo: windowTwo,
	}
}
