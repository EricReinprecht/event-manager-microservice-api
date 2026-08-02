package database

import (
	"gorm.io/gorm"

	"github.com/reinp/event-platform/backend/internal/models"
)

func Migrate(
	db *gorm.DB,
) error {

	if err := db.AutoMigrate(

		// auth
		&models.User{},
		&models.UserRole{},
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
	); err != nil {
		return err
	}

	return migrateLegacyUserRoles(db)
}

func migrateLegacyUserRoles(db *gorm.DB) error {
	return db.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_schema = current_schema()
				  AND table_name = 'users'
				  AND column_name = 'role'
			) THEN
				INSERT INTO user_roles (id, user_id, role, created_at)
				SELECT gen_random_uuid(), id, role, NOW()
				FROM users
				WHERE role IS NOT NULL
				ON CONFLICT (user_id, role) DO NOTHING;
			END IF;
		END $$;
	`).Error
}
