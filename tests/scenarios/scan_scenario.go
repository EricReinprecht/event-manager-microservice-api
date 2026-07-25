package scenarios

import (
	"testing"

	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"gorm.io/gorm"
)

type ScanScenario struct {
	DB *gorm.DB

	Staff models.User

	Customer models.User

	Party models.Party

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

	base := CreateTicketScenario(
		t,
		db,
		clock,
		requiresVerification,
	)

	return ScanScenario{
		DB: db,

		Staff: base.Staff,

		Customer: base.Customer,

		Party: base.Party,

		Category: base.Category,

		TicketCategory: base.TicketCategory,

		Window: base.Window,

		Ticket: base.Ticket,
	}
}
