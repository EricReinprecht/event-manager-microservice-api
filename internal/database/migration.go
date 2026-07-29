package database

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func Migrate(
	db *gorm.DB,
) error {

	return db.AutoMigrate(

		&models.User{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.EmailVerification{},

		&models.Category{},
		&models.Party{},
		&models.Media{},
		&models.PartyMedia{},

		&models.Ticket{},
		&models.TicketCategory{},
		&models.TicketAccessWindow{},
		&models.TicketScan{},

		&models.Purchase{},
		&models.PurchaseItem{},

		&models.PartyMember{},
		&models.PartyMemberRole{},

		&models.PaymentEvent{},
	)
}
