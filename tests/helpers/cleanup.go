package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func CleanDatabase(db *gorm.DB) error {

	tables := []interface{}{

		&models.TicketScan{},

		&models.PaymentEvent{},

		&models.PurchaseItem{},

		&models.Ticket{},

		&models.TicketAccessWindow{},

		&models.TicketCategory{},

		&models.Purchase{},

		&models.PartyMemberRole{},

		&models.PartyMember{},

		&models.Party{},

		&models.Category{},

		&models.User{},
	}

	for _, table := range tables {

		err := db.
			Unscoped().
			Session(
				&gorm.Session{
					AllowGlobalUpdate: true,
				},
			).
			Delete(table).
			Error

		if err != nil {
			return err
		}
	}

	return nil
}
