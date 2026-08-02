package database

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func Migrate(
	db *gorm.DB,
) error {

	return db.AutoMigrate(

		// auth
		&models.User{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.EmailVerification{},

		// media
		&models.Media{},

		// artist system
		&models.Artist{},

		// party system
		&models.PartyCategory{},
		&models.Party{},
		&models.PartyMedia{},

		// party permissions
		&models.PartyMember{},
		&models.PartyMemberRole{},

		// organizational staff
		&models.StaffMember{},

		// lineup system
		&models.PartyStage{},
		&models.ArtistSlot{},

		// ticketing
		&models.TicketCategory{},
		&models.TicketAccessWindow{},
		&models.Ticket{},
		&models.TicketScan{},

		// purchases
		&models.Purchase{},
		&models.PurchaseItem{},

		// payments
		&models.PaymentEvent{},
	)
}
