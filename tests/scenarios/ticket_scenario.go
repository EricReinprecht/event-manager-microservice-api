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

	Window *models.TicketAccessWindow

	Purchase models.Purchase

	Ticket models.Ticket
}

// =====================
// OPTIONS
// =====================

type TicketScenarioOption func(*TicketScenarioConfig)

type TicketScenarioConfig struct {
	CreateAccessWindow bool

	CreatePurchase bool
}

func DefaultTicketScenarioConfig() TicketScenarioConfig {

	return TicketScenarioConfig{

		CreateAccessWindow: true,

		CreatePurchase: true,
	}
}

func WithoutAccessWindow() TicketScenarioOption {

	return func(
		config *TicketScenarioConfig,
	) {

		config.CreateAccessWindow = false
	}
}

func WithoutPurchase() TicketScenarioOption {

	return func(
		config *TicketScenarioConfig,
	) {
		config.CreatePurchase = false
	}
}

// =====================
// SCENARIO CREATOR
// =====================

func CreateTicketScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
	requiresVerification bool,
	options ...TicketScenarioOption,
) TicketScenario {

	t.Helper()

	// APPLY CONFIG OPTIONS FIRST

	config := DefaultTicketScenarioConfig()

	for _, option := range options {

		option(&config)
	}

	// USER

	staff := fixtures.User()

	if err := db.Create(&staff).Error; err != nil {

		t.Fatal(err)
	}

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

	// PARTY ROLE

	AddPartyRole(
		t,
		db,
		staff.ID,
		party.ID,
		enum.RoleStaff,
	)

	// TICKET CATEGORY

	ticketCategory := models.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		Price: 100,

		Capacity: 100,

		PartyID: party.ID,

		RequiresVerification: requiresVerification,
	}

	if err := db.Create(&ticketCategory).Error; err != nil {

		t.Fatal(err)
	}

	// ACCESS WINDOW

	var window *models.TicketAccessWindow

	if config.CreateAccessWindow {

		createdWindow := models.TicketAccessWindow{

			ID: uuid.New(),

			TicketCategoryID: ticketCategory.ID,

			StartsAt: clock.Now().Add(-time.Hour),

			EndsAt: clock.Now().Add(time.Hour),
		}

		if err := db.Create(&createdWindow).Error; err != nil {

			t.Fatal(err)
		}

		window = &createdWindow
	}

	// PURCHASE

	var purchase models.Purchase

	if config.CreatePurchase {

		purchase = CreatePurchase(
			t,
			db,
			customer.ID,
			party.ID,
		)
	}

	// TICKET

	ticket := fixtures.Ticket()

	ticket.TicketCategoryID = ticketCategory.ID

	ticket.UserID = customer.ID

	if config.CreatePurchase {
		ticket.PurchaseID = purchase.ID
	}

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
