package helpers

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func CleanDatabase(db *gorm.DB) error {

	err := db.Exec(
		"SET CONSTRAINTS ALL DEFERRED",
	).Error

	if err != nil {
		return err
	}

	tables := []interface{}{

		&models.TicketScan{},

		&models.PaymentEvent{},

		&models.Ticket{},

		&models.PurchaseItem{},

		&models.Purchase{},

		&models.TicketAccessWindow{},

		&models.TicketCategory{},

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
