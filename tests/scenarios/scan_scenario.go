package scenarios

import (
	"testing"

	"github.com/google/uuid"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/service"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"gorm.io/gorm"
)

type ScanScenario struct {
	DB *gorm.DB

	TicketService *service.TicketService

	Staff    models.User
	Customer models.User

	Party    models.Party
	Category models.Category

	TicketCategory models.TicketCategory

	Window *models.TicketAccessWindow

	Ticket models.Ticket
}

func CreateScanScenario(
	t *testing.T,
	db *gorm.DB,
	clock *helpers.FakeClock,
	requiresVerification bool,
) ScanScenario {

	t.Helper()

	ticketService := helpers.NewTicketService(
		db,
		clock,
	)

	base := CreateTicketScenario(
		t,
		db,
		clock,
		requiresVerification,
	)

	return ScanScenario{

		DB: db,

		TicketService: ticketService,

		Staff: base.Staff,

		Customer: base.Customer,

		Party: base.Party,

		Category: base.Category,

		TicketCategory: base.TicketCategory,

		Window: base.Window,

		Ticket: base.Ticket,
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
