package scenarios

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

type TicketScenario struct {
	Staff models.User

	Customer models.User

	Party models.Party

	Category models.Category

	TicketCategory models.TicketCategory

	Window models.TicketAccessWindow

	Purchase models.Purchase

	Ticket models.Ticket
}

func CreateTicketScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
	requiresVerification bool,
) TicketScenario {

	t.Helper()

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}

	customer := fixtures.User()

	if err := db.Create(&customer).Error; err != nil {
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

		RequiresVerification: requiresVerification,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	window := models.TicketAccessWindow{
		ID: uuid.New(),

		TicketCategoryID: ticketCategory.ID,

		StartsAt: clock.Current.Add(-time.Hour),

		EndsAt: clock.Current.Add(time.Hour),
	}

	if err := db.Create(&window).Error; err != nil {
		t.Fatal(err)
	}

	purchase := CreatePurchase(
		t,
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

	return TicketScenario{
		Staff: staff,

		Customer: customer,

		Party: party,

		Category: category,

		TicketCategory: ticketCategory,

		Window: window,

		Purchase: purchase,

		Ticket: ticket,
	}
}
