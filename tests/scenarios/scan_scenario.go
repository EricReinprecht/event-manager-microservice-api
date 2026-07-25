package scenarios

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"gorm.io/gorm"
)

type ScanScenario struct {
	DB *gorm.DB

	TicketService *service.TicketService

	Staff models.User

	StaffTwo models.User

	Customer models.User

	Party models.Party

	Category models.Category

	TicketCategory models.TicketCategory

	Window models.TicketAccessWindow

	Ticket models.Ticket
}

func CreateScanScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
	options ...bool,
) ScanScenario {

	t.Helper()

	requiresVerification := false

	if len(options) > 0 {
		requiresVerification = options[0]
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

	// PARTY STAFF

	for _, userID := range []uuid.UUID{
		staffOne.ID,
		staffTwo.ID,
	} {

		AddPartyRole(
			t,
			db,
			userID,
			party.ID,
			enum.RoleStaff,
		)
	}

	// TICKET CATEGORY

	ticketCategory := models.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		RequiresVerification: requiresVerification,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	// ACCESS WINDOW

	window := models.TicketAccessWindow{

		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Now().Add(-time.Hour),

		EndsAt: clock.Now().Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	// PURCHASE

	purchase := CreatePurchase(
		t,
		db,
		customer.ID,
		party.ID,
	)

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = customer.ID

	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ScanScenario{

		DB: db,

		TicketService: ticketService,

		Staff: staffOne,

		StaffTwo: staffTwo,

		Customer: customer,

		Party: party,

		Category: category,

		TicketCategory: ticketCategory,

		Window: window,

		Ticket: ticket,
	}
}

func CreateAdditionalTicket(
	t *testing.T,
	db *gorm.DB,
	categoryID uuid.UUID,
	userID uuid.UUID,
	partyID uuid.UUID,
) models.Ticket {

	purchase := CreatePurchase(
		t,
		db,
		userID,
		partyID,
	)

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = categoryID

	ticket.UserID = userID

	ticket.PurchaseID = purchase.ID

	if err := db.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	return ticket
}

func CreatePendingTicketScan(
	t *testing.T,
	db *gorm.DB,
	ticket models.Ticket,
	window models.TicketAccessWindow,
	scanner models.User,
	at time.Time,
) models.TicketScan {

	t.Helper()

	scan := models.TicketScan{

		ID: uuid.New(),

		TicketID: ticket.ID,

		TicketAccessWindowID: window.ID,

		ScannedByID: scanner.ID,

		ScannedAt: at,

		Status: enum.TicketScanPending,
	}

	if err := db.Create(&scan).Error; err != nil {
		t.Fatal(err)
	}

	return scan
}
